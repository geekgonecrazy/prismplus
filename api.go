package main

import (
	"github.com/geekgonecrazy/prismplus/controllers"
	"github.com/geekgonecrazy/prismplus/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func apiServer() {
	sessions.InitializeSessionStore()

	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.CORS())

	keyAuthConfig := middleware.KeyAuthConfig{KeyLookup: "header:Authorization", Validator: validateAdminKey}

	router.GET("/api/v1/streamers", controllers.GetStreamersHandler, middleware.KeyAuthWithConfig(keyAuthConfig))
	router.POST("/api/v1/streamers", controllers.CreateStreamerHandler, middleware.KeyAuthWithConfig(keyAuthConfig))
	router.GET("/api/v1/streamers/:streamer", controllers.GetStreamerHandler, middleware.KeyAuthWithConfig(keyAuthConfig))
	router.DELETE("/api/v1/streamers/:streamer", controllers.DeleteStreamerHandler, middleware.KeyAuthWithConfig(keyAuthConfig))

	router.GET("/api/v1/streamer", controllers.GetMyStreamerHandler)
	router.GET("/api/v1/streamer/destinations", controllers.GetMyStreamerDestinationsHandler)
	router.POST("/api/v1/streamer/destinations", controllers.CreateMyStreamerDestinationHandler)
	router.DELETE("/api/v1/streamer/destinations/:destination", controllers.RemoveMyStreamerDestinationHandler)

	router.GET("/api/v1/sessions", controllers.GetSessionsHandler, middleware.KeyAuthWithConfig(keyAuthConfig))
	router.POST("/api/v1/sessions", controllers.CreateSessionHandler, middleware.KeyAuthWithConfig(keyAuthConfig))
	router.GET("/api/v1/sessions/:session", controllers.GetSessionHandler)
	router.POST("/api/v1/sessions/:session/destinations", controllers.AddDestinationHandler)
	router.GET("/api/v1/sessions/:session/destinations", controllers.GetDestinationsHandler)
	router.DELETE("/api/v1/sessions/:session/destinations/:destination", controllers.RemoveDestinationHandler)
	router.DELETE("/api/v1/sessions/:session", controllers.DeleteSessionHandler)

	router.Logger.Fatal(router.Start(":5383"))
}

func validateAdminKey(key string, c echo.Context) (bool, error) {
	if key == *adminKey {
		return true, nil
	}

	return false, nil
}
