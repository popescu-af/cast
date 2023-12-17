package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"app/logic"
	"logger"
	"source/directory"
)

var lg *logic.Logic

var (
	deviceScanTimeout = 1000 * time.Millisecond
	supportedDevices  = []string{"chromecast"}
)

func init() {
	var err error
	lg, err = logic.New("")
	if err != nil {
		panic(err)
	}

	mediaLibraryPath := os.Getenv("MEDIA_LIBRARY_PATH")
	if mediaLibraryPath == "" {
		panic("MEDIA_LIBRARY_PATH not set")
	}

	lib, err := directory.New(mediaLibraryPath, logic.SupportedExtensions)
	if err != nil {
		panic(err)
	}
	lg.Manager.AddLibrary(lib)

	messages := lg.FindDevices(deviceScanTimeout, supportedDevices)
	for _, m := range messages {
		logger.Log.Printf(m)
	}
}

func main() {
	// service
	r := gin.Default()
	// CORS
	r.Use(cors.Default())
	// routes
	r.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	r.GET("/media", getMedia)
	r.GET("/devices", getDevices)
	r.POST("/load", loadMedia)
	r.POST("/playback", controlPlayback)
	r.GET("/playback", getPlaybackInfo)
	r.GET("/logs", getLogs)
	// run
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
