package main

import (
	"cloud-run-sample/eth"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// GitCommitLog is set at build-time
var GitCommitLog string

func init() {
	fmt.Printf("GIT log: %s\n", GitCommitLog)
}

const (
	NftContract = "0xfF646D99fB94bb20439429c8fe0EE2F58090FA14"
	BscRpc      = "https://data-seed-prebsc-1-s1.binance.org:8545"
)

type (
	Data struct {
		Data string `json:"data"`
	}

	CheckIn struct {
		Address   string `json:"address"`
		TokenId   string `json:"token_id"`
		Signature string `json:"signature"`
		Contract  string `json:"contract"`
		Timestamp int    `json:"timestamp"`
	}
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.GET("/", hello)
	e.GET("/git", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("GIT log: %s\n", GitCommitLog))
	})
	e.GET("/time", timeNow)
	e.POST("/check-in", checkIn)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func hello(c echo.Context) error {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	return c.String(http.StatusOK, fmt.Sprintf("Hello %s!\n", name))
}

func timeNow(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("Now %s!\n", time.Now().String()))
}

func checkIn(c echo.Context) error {
	client, err := ethclient.Dial(BscRpc)
	if err != nil {
		return err
	}
	var data Data
	if err := c.Bind(&data); err != nil {
		return err
	}
	//address::token_id::signature::timestamp
	checkInSplit := strings.Split(data.Data, "::")
	if len(checkInSplit) != 4 || len(checkInSplit[0]) != 42 || len(checkInSplit[3]) != 10 || len(checkInSplit[2]) != 132 {
		return c.String(http.StatusBadRequest, "invalid data. format: address::token_id::signature::timestamp")
	}
	timestamp, err := strconv.Atoi(checkInSplit[3])
	if err != nil {
		return err
	}
	if timestamp < int(time.Now().Unix())-60 {
		return c.String(http.StatusBadRequest, "expired")
	}

	if !eth.VerifySig(checkInSplit[0], checkInSplit[2], []byte(fmt.Sprintf("%s::%s", checkInSplit[1], checkInSplit[3]))) {
		return c.String(http.StatusBadRequest, "invalid signature")
	}

	tokenId, err := strconv.Atoi(checkInSplit[1])
	nft, err := eth.NewNft(common.HexToAddress(NftContract), client)
	if err != nil {
		return err
	}
	ownerOf, err := nft.OwnerOf(nil, big.NewInt(int64(tokenId)))
	if err != nil {
		return err
	}
	if common.HexToAddress(checkInSplit[0]) != ownerOf {
		return c.String(http.StatusBadRequest, "invalid owner")
	}

	checkInData := CheckIn{
		Address:   checkInSplit[0],
		TokenId:   checkInSplit[1],
		Signature: checkInSplit[2],
		Contract:  NftContract,
		Timestamp: timestamp,
	}
	return c.JSON(http.StatusOK, checkInData)
}
