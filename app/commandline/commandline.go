package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"app/logic"
	"model/session"
	"source/directory"
)

var lg *logic.Logic

func init() {
	var err error
	lg, err = logic.New("")
	if err != nil {
		panic(err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		params := strings.Split(strings.ReplaceAll(input, "\n", ""), " ")
		cmd := params[0]
		params = params[1:]

		switch cmd {
		case "q":
			fmt.Println("quitting")
			return
		case "h":
			fmt.Println(help)
		case "s":
			err = scanForDevices(params)
		case "ld":
			err = listDevices()
		case "ll":
			err = listLibraries()
		case "lm":
			err = listMedia(params)
		case "ad":
			err = addDirectoryLibrary(params)
		case "load":
			err = loadMedia(params)
		case "ctrl":
			err = commandMedia(params)
		default:
			fmt.Println("unknown command")
		}
		if err != nil {
			fmt.Println("error -", err)
		}
	}
}

var help = `Commands:
  q                                    quit
  h                                    show the help
  s                                    scan for devices
  ld                                   list devices
  ll                                   list libraries
  lm library_id                        list media
  ad dir_path                          add library
  load device_id media_id subtitle_id  load media on device
  ctrl cmd [params...]                 issue media control command

Session commands:
  play [from_beginning]  play media
  pause                  pause playing
  seek absolute_per      seek video at given absolute percentage
  fwd seconds            seek video forward by given seconds
  rev seconds            seek video backward by given seconds
`

var errWrongNumberOfParams = errors.New("wrong number of parameters for command")

func scanForDevices(params []string) error {
	if len(params) > 1 {
		return errWrongNumberOfParams
	}
	var timeout = 1000 * time.Millisecond
	if len(params) == 1 {
		var err error
		timeout, err = time.ParseDuration(params[0])
		if err != nil {
			return err
		}
	}
	fmt.Println("scanning for devices...")
	messages := lg.FindDevices(timeout, nil)
	for _, message := range messages {
		fmt.Println(message)
	}
	return nil
}

func listDevices() error {
	for _, device := range lg.Manager.ListDevices() {
		fmt.Println(device.ID())
	}
	return nil
}

func listLibraries() error {
	for _, lib := range lg.Manager.ListLibraries() {
		fmt.Println(lib.ID())
	}
	return nil
}

func listMedia(params []string) error {
	if len(params) != 1 {
		return errWrongNumberOfParams
	}
	media, err := lg.Manager.ListMedia(params[0])
	if err != nil {
		return err
	}
	for _, m := range media {
		fmt.Println(m.ID(), m.URI())
	}
	return nil
}

func addDirectoryLibrary(params []string) error {
	if len(params) != 1 {
		return errWrongNumberOfParams
	}
	lib, err := directory.New(params[0], logic.SupportedExtensions)
	if err != nil {
		return err
	}
	lg.Manager.AddLibrary(lib)
	return nil
}

var currentSession session.Session

func loadMedia(params []string) error {
	if len(params) < 2 {
		return errWrongNumberOfParams
	}
	var subtitleID string
	if len(params) == 3 {
		subtitleID = params[2]
	}
	var err error
	currentSession, err = lg.Manager.CreateSession(params[0], params[1], subtitleID)
	if err != nil {
		return err
	}
	fmt.Println(currentSession.ID())
	return nil
}

func commandMedia(params []string) error {
	if currentSession == nil {
		return errors.New("no media loaded")
	}
	if len(params) < 1 {
		return errWrongNumberOfParams
	}
	cmd := params[0]
	params = params[1:]
	switch cmd {
	case "play":
		var fromBeginning bool
		if len(params) > 1 {
			return errWrongNumberOfParams
		}
		if len(params) == 1 {
			var err error
			fromBeginning, err = strconv.ParseBool(params[0])
			if err != nil {
				return err
			}
		}
		return currentSession.Play(fromBeginning)
	case "pause":
		return currentSession.Pause()
	case "seek":
		if len(params) != 1 {
			return errWrongNumberOfParams
		}
		percentage, err := strconv.ParseFloat(params[0], 64)
		if err != nil {
			return err
		}
		return currentSession.Seek(percentage)
	case "fwd":
		if len(params) != 1 {
			return errWrongNumberOfParams
		}
		seconds, err := strconv.ParseFloat(params[0], 64)
		if err != nil {
			return err
		}
		return currentSession.Fwd(seconds)
	case "rev":
		if len(params) != 1 {
			return errWrongNumberOfParams
		}
		seconds, err := strconv.ParseFloat(params[0], 64)
		if err != nil {
			return err
		}
		return currentSession.Rev(seconds)
	}
	return errors.New("unsupported player command")
}
