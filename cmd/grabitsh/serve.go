package grabitsh

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Run:   startWebServer,
}

func startWebServer(cmd *cobra.Command, args []string) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to Grabit.sh!")
	})
	e.Logger.Fatal(e.Start(":42069"))
}
