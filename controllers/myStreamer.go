package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/geekgonecrazy/prismplus/models"
	"github.com/geekgonecrazy/prismplus/sessions"
	"github.com/geekgonecrazy/prismplus/store"
	"github.com/geekgonecrazy/prismplus/streamers"
	"github.com/labstack/echo/v4"
)

func getStreamKeyFromHeader(c echo.Context) (string, bool) {
	authorizationHeader := c.Request().Header.Get("Authorization")

	split := strings.Split(authorizationHeader, " ")

	if len(split) < 2 {
		return "", false
	}

	if len(split[1]) == 0 {
		return "", false
	}

	return split[1], true
}

func GetMyStreamerHandler(c echo.Context) error {
	streamKey, ok := getStreamKeyFromHeader(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	streamer, err := streamers.GetStreamerByStreamKey(streamKey)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	myStreamer := models.MyStreamer{
		Streamer: streamer,
	}

	session, _ := sessions.GetSession(streamKey)

	if session != nil && session.Running {
		myStreamer.Live = true
	}

	return c.JSON(http.StatusOK, myStreamer)
}

func CreateMyStreamerDestinationHandler(c echo.Context) error {
	streamKey, ok := getStreamKeyFromHeader(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	myStreamer, err := streamers.GetStreamerByStreamKey(streamKey)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	destinationPayload := models.Destination{}

	if err := c.Bind(&destinationPayload); err != nil {
		return err
	}

	if err := streamers.AddDestination(myStreamer, destinationPayload); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func RemoveMyStreamerDestinationHandler(c echo.Context) error {
	streamKey, ok := getStreamKeyFromHeader(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	myStreamer, err := streamers.GetStreamerByStreamKey(streamKey)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	destination := c.Param("destination")
	id, err := strconv.Atoi(destination)
	if err != nil {
		return c.String(http.StatusBadRequest, "Not Found")
	}

	if err := streamers.RemoveDestination(myStreamer, id); err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusAccepted)
}

func GetMyStreamerDestinationsHandler(c echo.Context) error {
	streamKey, ok := getStreamKeyFromHeader(c)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	myStreamer, err := streamers.GetStreamerByStreamKey(streamKey)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, myStreamer.Destinations)
}
