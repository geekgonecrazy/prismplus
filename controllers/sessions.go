package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/geekgonecrazy/prismplus/models"
	"github.com/geekgonecrazy/prismplus/sessions"
	"github.com/labstack/echo/v4"
)

func GetSessionsHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, sessions.GetSessions())
}

func CreateSessionHandler(c echo.Context) error {
	sessionPayload := models.SessionPayload{}

	if err := c.Bind(&sessionPayload); err != nil {
		return err
	}

	if err := sessions.CreateSession(sessionPayload); err != nil {
		if err.Error() == "Already Exists" {
			return c.NoContent(http.StatusConflict)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, sessionPayload)
}

func GetSessionHandler(c echo.Context) error {
	key := c.Param("session")

	session, err := sessions.GetSession(key)
	if err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, session)
}

func GetDestinationsHandler(c echo.Context) error {
	key := c.Param("session")

	session, err := sessions.GetSession(key)
	if err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	destinations := session.GetDestinations()

	return c.JSON(http.StatusOK, destinations)
}

func AddDestinationHandler(c echo.Context) error {
	key := c.Param("session")

	session, err := sessions.GetSession(key)
	if err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	destinationPayload := models.Destination{}

	if err := c.Bind(&destinationPayload); err != nil {
		return err
	}

	if err := session.AddDestination(destinationPayload); err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func RemoveDestinationHandler(c echo.Context) error {
	key := c.Param("session")
	destination := c.Param("destination")

	id, err := strconv.Atoi(destination)
	if err != nil {
		return c.String(http.StatusBadRequest, "Not Found")
	}

	session, err := sessions.GetSession(key)
	if err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	if err := session.RemoveDestination(id); err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusAccepted)
}

func DeleteSessionHandler(c echo.Context) error {
	key := c.Param("session")

	session, err := sessions.GetSession(key)
	if err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusInternalServerError)
	}

	session.EndSession()

	return c.NoContent(http.StatusAccepted)
}
