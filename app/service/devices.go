package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getDevices(c *gin.Context) {
	rescanParam := c.DefaultQuery("rescan", "false")
	var rescan bool
	var err error
	rescan, err = strconv.ParseBool(rescanParam)
	if err != nil {
		rescan = false
	}
	var messages []string
	if rescan {
		messages = lg.FindDevices(deviceScanTimeout, supportedDevices)
	}
	devices := lg.Manager.ListDevices()
	devicesJSON := make([]deviceJSON, 0)
	for _, d := range devices {
		devicesJSON = append(devicesJSON, deviceJSON{
			ID: d.ID(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data":     devicesJSON,
		"messages": messages,
	})
}

type deviceJSON struct {
	ID string `json:"id"`
}
