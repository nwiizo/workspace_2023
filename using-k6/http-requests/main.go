package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/login", loginHandler)

	e.Start(":8080")
}

func loginHandler(c echo.Context) error {
	creds := new(Credentials)
	if err := json.NewDecoder(c.Request().Body).Decode(creds); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if creds.Email == "aaa" && creds.Password == "bbb" {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Authenticated",
		})
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{
		"error": "Invalid email or password",
	})
}
