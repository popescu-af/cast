package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getMedia(c *gin.Context) {
	libs := lg.Manager.ListLibraries()
	allMedia := make([]mediaJSON, 0)
	for _, lib := range libs {
		media, err := lib.List()
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"error": err.Error(),
				},
			)
		}
		for _, m := range media {
			allMedia = append(allMedia, mediaJSON{
				ID:  m.ID(),
				URI: m.URI(),
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": allMedia,
	})
}

type mediaJSON struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}
