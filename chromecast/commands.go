package chromecast

import (
	"fmt"
	"model/library"
	"path/filepath"

	"github.com/looplab/fsm"
)

func (s *state) initiateConnection(e *fsm.Event) {
	requests := []*pkg{
		newConnectRequest(s.senderID, s.receiverID, "CONNECT"),
		newGetReceiverStatusRequest(s.senderID, s.receiverID, s.nextRequestID()),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) openSession(e *fsm.Event) {
	requests := []*pkg{
		newLaunchRequest(s.senderID, s.receiverID, s.nextRequestID(), s.appID),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) cacheMedia(e *fsm.Event) {
	if len(e.Args) == 0 {
		e.Err = fmt.Errorf("no media provided")
		return
	}
	s.mediaPlaying = e.Args[0].(library.Media)
	if len(e.Args) == 2 && e.Args[1] != nil {
		subtitle := e.Args[1].(library.Media)
		extension := filepath.Ext(subtitle.URI())
		if extension != ".vtt" {
			vtt, err := s.subtitleFS.convertToVTT(subtitle)
			if err != nil {
				e.Err = err
				return
			}
			s.subtitle = vtt
		} else {
			s.subtitle = subtitle
		}
	}
}

func (s *state) loadCachedMedia(e *fsm.Event) {
	if s.mediaPlaying == nil {
		e.Err = fmt.Errorf("media is missing")
		return
	}
	var uriVideo = s.mediaPlaying.URI()
	var uriSubtitle string
	if s.subtitle != nil {
		uriSubtitle = s.subtitle.URI()
	}
	requests := []*pkg{
		newConnectRequest(s.senderID, s.sessionID, "CONNECT"),
		newLoadRequest(s.senderID, s.sessionID, uriVideo, uriSubtitle, s.nextRequestID()),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) play(e *fsm.Event) {
	requests := []*pkg{
		newMediaSessionRequest(s.senderID, s.sessionID, "PLAY", s.nextRequestID(), s.mediaSessionID),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) pause(e *fsm.Event) {
	requests := []*pkg{
		newMediaSessionRequest(s.senderID, s.sessionID, "PAUSE", s.nextRequestID(), s.mediaSessionID),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) stop(e *fsm.Event) {
	requests := []*pkg{
		newMediaSessionRequest(s.senderID, s.sessionID, "STOP", s.nextRequestID(), s.mediaSessionID),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) closeSession(e *fsm.Event) {
	requests := []*pkg{
		newConnectRequest(s.senderID, s.sessionID, "CLOSE"),
		newConnectRequest(s.senderID, s.receiverID, "CLOSE"),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) seek(e *fsm.Event) {
	if s.playbackInfo.Duration == 0 {
		e.Err = fmt.Errorf("cannot seek, duration is zero")
		return
	}
	var t = e.Args[1].(seekType)
	var amountParam = e.Args[2].(float64)
	var seekOffset float64
	switch t {
	case seekTypeRelative:
		seekOffset = s.playbackInfo.Position + amountParam
		if seekOffset < 0 {
			seekOffset = 0
		} else if seekOffset > s.playbackInfo.Duration {
			seekOffset = s.playbackInfo.Duration
		}
	case seekTypeAbsolutePercentage:
		if amountParam < 0 {
			amountParam = 0
		} else if amountParam > 100 {
			amountParam = 100
		}
		seekOffset = amountParam * s.playbackInfo.Duration / 100
	}
	requests := []*pkg{
		newSeekMediaSessionRequest(s.senderID, s.sessionID, s.nextRequestID(), s.mediaSessionID, seekOffset),
	}
	e.Err = s.sendRequests(requests)
}

func (s *state) requestMediaStatus(e *fsm.Event) {
	requests := []*pkg{
		newGetMediaStatusRequest(s.senderID, s.sessionID, s.mediaSessionID, s.nextRequestID()),
	}
	e.Err = s.sendRequests(requests)
}
