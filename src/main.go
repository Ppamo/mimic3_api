package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"ppamo/api/common"
	"ppamo/api/config"
	"ppamo/api/convertion"
	"ppamo/api/handlers"

	"github.com/labstack/echo/v4"
)

const (
	PORT = 8080
)

var (
	Default501Error = common.DefaultError{Status: 501, Description: "Internal Server Error"}
)

func main() {
	config_path := os.Getenv("CONFIG_PATH")
	if _, err := os.Stat(config_path); err != nil {
		log.Fatalf("++> Error loading config:\n%v", err)
	}
	config.LoadConfig(os.Getenv("CONFIG_PATH"))
	e := echo.New()
	conv, err := convertion.NewConverter()
	if err != nil {
		log.Fatalf("++> Error creating new converter:\n%v", err)
	}
	hand, err := handlers.NewHandler(&conv)
	if err != nil {
		log.Fatalf("++> Error creating new handler:\n%v", err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, common.DefaultError{
			Status:      500,
			Description: "Method Not Yet Implemented"},
		)
	})
	e.GET("/profiles", func(c echo.Context) error {
		res, err := hand.GetProfiles()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Default501Error)
		}
		return c.JSON(http.StatusOK, res)
	})
	e.GET("/effects", func(c echo.Context) error {
		res, err := hand.GetEffects()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Default501Error)
		}
		return c.JSON(http.StatusOK, res)
	})
	e.POST("/convert", func(c echo.Context) error {
		req := common.ConvertRequest{}
		err := json.NewDecoder(c.Request().Body).Decode(&req)
		if err != nil {
			log.Printf("++> ERROR: Failed to unmarshal", err)
			return c.JSON(http.StatusBadRequest, common.DefaultError{
				Status:      400,
				Description: fmt.Sprintf("Bad Request:\n%s", err),
			})
		}
		log.Printf("++> Converting request:\n%v", req)
		res, err := hand.Convert(&req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, res)
		}
		return c.JSON(http.StatusOK, res)
	})
	log.Printf("++> Starting server at port %d", PORT)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", PORT)))
}
