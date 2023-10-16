package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ppamo/api/convertion"

	"github.com/labstack/echo/v4"
)

const (
	PORT = 8080
)

type ResponseStuct struct {
	Status      int    `json:"status"`
	Body        []byte `json:"body,omitempty"`
	Description string `json:"description"`
}

type RequestStruct struct {
	Text    string `json:"text"`
	Profile string `json:"profile"`
}

func main() {
	e := echo.New()
	con, err := convertion.NewConverter()
	if err != nil {
		log.Fatalf("--> Error creating converter:\n%v", err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, ResponseStuct{Status: 200, Description: "OK"})
	})
	e.POST("/convert", func(c echo.Context) error {
		req := RequestStruct{}
		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			return c.JSON(http.StatusNotImplemented, ResponseStuct{Status: 400, Description: "NOK"})
		}
		_, err := con.Convert(&convertion.ConvertionRequest{
			Text:    req.Text,
			Profile: req.Profile,
		})
		if err != nil {
			return c.JSON(http.StatusNotImplemented, ResponseStuct{Status: 400, Description: "NOK"})
		}
		return c.JSON(http.StatusNotImplemented, ResponseStuct{Status: 200, Description: "OK"})
	})
	log.Printf("Starting server at port %d", PORT)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", PORT)))
}
