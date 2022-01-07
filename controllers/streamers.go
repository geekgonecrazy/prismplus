package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/geekgonecrazy/prismplus/models"
	"github.com/geekgonecrazy/prismplus/store"
	"github.com/geekgonecrazy/prismplus/streamers"
	"github.com/labstack/echo/v4"
)

func GetStreamersHandler(c echo.Context) error {
	s, err := streamers.GetStreamers()
	if err != nil {
		log.Println("Error:", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, s)
}

func CreateStreamerHandler(c echo.Context) error {
	streamerPayload := models.StreamerCreatePayload{}

	if err := c.Bind(&streamerPayload); err != nil {
		return err
	}

	if streamerPayload.Name == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	streamer, err := streamers.CreateStreamer(streamerPayload)
	if err != nil {
		if err.Error() == "Already Exists" {
			return c.NoContent(http.StatusConflict)
		}

		log.Println("Error:", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, streamer)
}

func GetStreamerHandler(c echo.Context) error {
	key := c.Param("streamer")

	id, err := strconv.Atoi(key)
	if err != nil {
		return c.String(http.StatusBadRequest, "Not Found")
	}

	streamer, err := streamers.GetStreamer(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, streamer)
}

func DeleteStreamerHandler(c echo.Context) error {
	key := c.Param("streamer")

	id, err := strconv.Atoi(key)
	if err != nil {
		return c.String(http.StatusBadRequest, "Not Found")
	}

	if err := streamers.DeleteStreamer(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusAccepted)
}
