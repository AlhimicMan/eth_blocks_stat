package http_server

import (
	"eth_blocks_stat/core/activity_calculator"
	"eth_blocks_stat/infrastructure/adapters"
	"eth_blocks_stat/infrastructure/config"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	cfg config.ServiceConfig
}

func RunServer(cfg config.ServiceConfig) error {
	router := echo.New()
	listenAddress := fmt.Sprintf("%s:%d", cfg.HTTPServerHost, cfg.HTTPServerPort)
	httpServer := &http.Server{
		Addr:    listenAddress,
		Handler: router,
	}
	fmt.Printf("Listening on %s\n", listenAddress)

	apiGroup := router.Group("/api")

	mountApiViews(apiGroup, cfg)

	for _, route := range router.Routes() {
		fmt.Printf("registered route %s:%s\n", route.Method, route.Path)
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("cannot run http server: %w", err)
	}
	return nil
}

func mountApiViews(apiGroup *echo.Group, cfg config.ServiceConfig) {
	apiGroup.GET("/top_active", func(c echo.Context) error {
		ctx := c.Request().Context()
		gBClient := adapters.NewGetBlockClient(cfg.APIKey)
		statModule := activity_calculator.NewCalculator(gBClient)
		res, err := statModule.RetrieveTopAddresses(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, res)
	})
}
