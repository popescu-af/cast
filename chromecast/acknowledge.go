package chromecast

import (
	"logger"
	"model/device"

	"github.com/looplab/fsm"
)

func (s *state) ackSession(e *fsm.Event) {
	receiverStatus := e.Args[0].(*receiverStatus)
	for _, app := range receiverStatus.Status.Applications {
		if app.AppID == s.appID {
			if app.SessionID == "" {
				logger.Log.Printf("ERROR: session id is empty")
			}
			s.sessionID = app.SessionID
			break
		}
	}
}

func (s *state) ackMediaSession(e *fsm.Event) {
	mediaStatus := e.Args[0].(*mediaStatus)
	for _, status := range mediaStatus.Status {
		if status.PlayerState == "IDLE" {
			if status.ExtendedStatus.PlayerState == "LOADING" {
				s.mediaSessionID = status.ExtendedStatus.MediaSessionId
				break
			}
		} else {
			s.mediaSessionID = status.MediaSessionID
			break
		}
	}
}

func (s *state) ackMediaStopped() {
	s.mediaSessionID = 0
	s.mediaPlaying = nil
	s.subtitle = nil
	s.playbackInfo = device.PlaybackInfo{}
}
