package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"model/session"
)

var currentSession session.Session

func loadMedia(c *gin.Context) {
	mediaID := c.DefaultQuery("mediaId", "")
	subtitleID := c.DefaultQuery("subtitleId", "")
	deviceID := c.DefaultQuery("deviceId", "")
	if mediaID == "" || deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing mediaId or deviceId",
		})
		return
	}
	// TODO: stop current session
	// if currentSession != nil {
	// 	currentSession.Stop()
	// }
	var err error
	currentSession, err = lg.Manager.CreateSession(deviceID, mediaID, subtitleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func controlPlayback(c *gin.Context) {
	if currentSession == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no playback in progress",
		})
		return
	}
	commandParam := c.DefaultQuery("command", "")
	if commandParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing command",
		})
		return
	}
	switch commandParam {
	case "play":
		restartParam := c.DefaultQuery("restart", "false")
		var restart bool
		var err error
		restart, err = strconv.ParseBool(restartParam)
		if err != nil {
			restart = false
		}
		err = currentSession.Play(restart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	case "pause":
		err := currentSession.Pause()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	case "stop":
		err := currentSession.Stop()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	case "seek":
		amountParam := c.DefaultQuery("amount", "") // 0-100
		var amount float64
		var err error
		amount, err = strconv.ParseFloat(amountParam, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid amount",
			})
			return
		}
		err = currentSession.Seek(amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	case "fwd", "rev":
		amountParam := c.DefaultQuery("amount", "")
		var amount float64
		var err error
		amount, err = strconv.ParseFloat(amountParam, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid amount",
			})
			return
		}
		if commandParam == "rev" {
			err = currentSession.Rev(amount * -1)
		} else {
			err = currentSession.Fwd(amount)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unsupported player command",
		})
	}
}

func getPlaybackInfo(c *gin.Context) {
	if currentSession == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no playback in progress",
		})
		return
	}
	playbackInfo, err := currentSession.GetPlaybackInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": playbackInfoJSON{
			Duration:   playbackInfo.Duration,
			Position:   playbackInfo.Position,
			PositionTS: playbackInfo.PositionTS,
			Playing:    playbackInfo.Playing,
		},
	})
}

type playbackInfoJSON struct {
	Duration   float64   `json:"duration"`
	Position   float64   `json:"position"`
	PositionTS time.Time `json:"positionTs"`
	Playing    bool      `json:"playing"`
}
