package chromecast

import (
	"context"
	"encoding/json"
	"fmt"
	"logger"
	"time"

	"github.com/looplab/fsm"

	"model/device"
	"model/library"
)

type seekType int

const (
	seekTypeRelative seekType = iota
	seekTypeAbsolutePercentage
)

type state struct {
	ctx            context.Context
	transport      *transport
	fsm            *fsm.FSM
	commandSync    *commandSync
	subtitleFS     *subtitleFS
	senderID       string
	receiverID     string
	appID          string
	reqIDGen       int
	sessionID      string
	mediaSessionID int
	mediaPlaying   library.Media
	subtitle       library.Media
	playbackInfo   device.PlaybackInfo
}

func newState(ctx context.Context, t *transport) (*state, error) {
	s := &state{
		ctx:         ctx,
		transport:   t,
		commandSync: newCommandSync(),
		senderID:    "sender-0",
		receiverID:  "receiver-0",
		appID:       "CC1AD845", // Default Media Receiver
	}
	var err error
	if s.subtitleFS, err = newSubtitleFS(); err != nil {
		return nil, err
	}
	s.createInternalFSM()
	return s, nil
}

func (s *state) close() error {
	return s.subtitleFS.fileserver.Close()
}

func (s *state) createInternalFSM() {
	s.fsm = fsm.NewFSM(
		"not_casting",
		fsm.Events{
			{
				Name: "request_cast",
				Src: []string{
					"casting",
					"not_casting",
					"loading",
					"reconnecting",
					"connecting",
					"opening_session",
					"controlling",
					"stopping",
				},
				Dst: "connecting",
			},

			{Name: "receiver_close", Src: []string{"loading"}, Dst: "reconnecting"},
			{Name: "receiver_status", Src: []string{"reconnecting"}, Dst: "opening_session"},

			{Name: "receiver_status", Src: []string{"connecting"}, Dst: "opening_session"},
			{Name: "receiver_status", Src: []string{"opening_session"}, Dst: "loading"},
			{Name: "media_status", Src: []string{"loading"}, Dst: "casting"},

			{Name: "request_control", Src: []string{"casting"}, Dst: "controlling"},
			{Name: "media_status", Src: []string{"controlling"}, Dst: "casting"},

			{Name: "request_stop", Src: []string{"casting"}, Dst: "stopping"},
			{Name: "media_status", Src: []string{"stopping"}, Dst: "not_casting"},
		},
		fsm.Callbacks{
			"enter_state": s.enterState,
		},
	)
}

func (s *state) emitEvent(event string, args ...interface{}) error {
	return s.fsm.Event(s.ctx, event, args...)
}

func (s *state) enterState(ctx context.Context, e *fsm.Event) {
	logger.Log.Printf("FSM state change %s --> %s", e.Src, e.Dst)
	switch e.Dst {
	case "not_casting":
		s.ackMediaStopped()
		s.closeSession(e)
		s.commandSync.notifyDone()
	case "casting":
		s.ackMediaSession(e)
		// TODO: check if really casting, fallback to not_casting if not
		s.commandSync.notifyDone()
	case "connecting":
		s.cacheMedia(e)
		s.initiateConnection(e)
	case "reconnecting":
		s.initiateConnection(e)
	case "opening_session":
		s.openSession(e)
	case "loading":
		s.ackSession(e)
		s.loadCachedMedia(e)
	case "controlling":
		command := e.Args[0].(string)
		switch command {
		case "PLAY":
			s.play(e)
		case "PAUSE":
			s.pause(e)
		case "SEEK":
			s.seek(e)
		case "GET_STATUS":
			s.requestMediaStatus(e)
		}
	case "stopping":
		s.stop(e)
	}
}

func (s *state) sendRequests(requests []*pkg) error {
	if s.transport == nil {
		return fmt.Errorf("transport not set")
	}
	for _, r := range requests {
		if err := s.transport.send(r); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) handleReceivedMessage(message *pkg) error {
	var typeWrapper requestWithType
	if err := json.Unmarshal(message.payload, &typeWrapper); err != nil {
		return err
	}

	switch typeWrapper.Type {
	case "RECEIVER_STATUS":
		var status receiverStatus
		if err := json.Unmarshal(message.payload, &status); err != nil {
			return err
		}
		return s.emitEvent("receiver_status", &status)

	case "MEDIA_STATUS":
		var status mediaStatus
		if err := json.Unmarshal(message.payload, &status); err != nil {
			return err
		}
		if len(status.Status) == 0 {
			logger.Log.Printf("no status in media status response")
			return s.emitEvent("media_status", &status)
		}
		firstStatus := status.Status[0]
		playerState := firstStatus.PlayerState
		if firstStatus.Media.Duration > 0 || playerState == "IDLE" {
			s.playbackInfo.Duration = firstStatus.Media.Duration
		}
		s.playbackInfo.Position = firstStatus.CurrentTime
		s.playbackInfo.PositionTS = time.Now()
		s.playbackInfo.Playing = playerState == "PLAYING"
		return s.emitEvent("media_status", &status)

	case "LOAD_FAILED":
		return fmt.Errorf("not implemented")
	case "PING":
		requests := []*pkg{
			newPongMessage(),
		}
		return s.sendRequests(requests)
	case "CLOSE":
		if message.source == s.sessionID {
			return s.emitEvent("receiver_close")
		}
		return nil

	default:
		return fmt.Errorf("unknown message type: %s", typeWrapper.Type)
	}
}

func (s *state) nextRequestID() int {
	s.reqIDGen++
	return s.reqIDGen
}
