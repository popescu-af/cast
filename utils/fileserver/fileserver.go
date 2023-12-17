package fileserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type Instance struct {
	ServingURL string
	server     *http.Server
}

func New(path string) (*Instance, error) {
	ip, err := getOutboundIP()
	if err != nil {
		return nil, err
	}
	fs := http.FileServer(http.Dir(path))
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// for some reason subtitles are not showing without these headers set, videos do work without them
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin,access-control-allow-headers")
			fs.ServeHTTP(w, r)
		}),
	}
	go server.Serve(listener)

	port := listener.Addr().(*net.TCPAddr).Port
	return &Instance{
		server:     server,
		ServingURL: fmt.Sprintf("http://%s:%d", ip, port),
	}, nil
}

func (i *Instance) Close() error {
	return i.server.Shutdown(context.Background())
}

// Get preferred outbound ip of this machine
// from https://stackoverflow.com/a/37382208
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP, nil
}
