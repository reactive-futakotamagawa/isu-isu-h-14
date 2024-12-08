package a

import ( // want "change import"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func initializeHandler(c echo.Context) error { // want "initialize"
	return c.String(http.StatusOK, "Hello, World!")
}

func b() {
	fmt.Println("Hello, World!")

	e := echo.New()
	e.POST("/initialize", initializeHandler)

	e.POST("/initialize", func(c echo.Context) error { // want "initialize"
		return c.String(http.StatusOK, "Hello, World!")
	})
}
