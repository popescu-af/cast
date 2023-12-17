package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"logger"
	"source/directory"
)

var logLibrary *directory.Library

func init() {
	// serve log files
	var err error
	logLibrary, err = directory.New(logger.Directory(), []string{".log"})
	if err != nil {
		panic(err)
	}
}

func getLogs(c *gin.Context) {
	logs, err := logLibrary.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	logPaths := make([]string, 0)
	for _, l := range logs {
		logPaths = append(logPaths, l.URI())
	}
	c.JSON(http.StatusOK, gin.H{
		"data": logPaths,
	})
}
