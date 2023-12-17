package chromecast

import "encoding/json"

type requestWithType struct {
	Type string `json:"type"`
}

func newPongMessage() *pkg {
	b, _ := json.Marshal(&requestWithType{Type: "PONG"})
	return &pkg{
		source:      "Tr@n$p0rt-0",
		destination: "Tr@n$p0rt-0",
		namespace:   "urn:x-cast:com.google.cast.tp.heartbeat",
		payload:     b,
	}
}

func newConnectRequest(s, d, t string) *pkg {
	b, _ := json.Marshal(&requestWithType{Type: t})
	return &pkg{
		source:      s,
		destination: d,
		namespace:   "urn:x-cast:com.google.cast.tp.connection",
		payload:     b,
	}
}

type requestWithTypeAndID struct {
	Type      string `json:"type"`
	RequestID int    `json:"requestId"`
}

func newGetReceiverStatusRequest(s, d string, id int) *pkg {
	b, _ := json.Marshal(&requestWithTypeAndID{Type: "GET_STATUS", RequestID: id})
	return &pkg{
		source:      s,
		destination: d,
		namespace:   "urn:x-cast:com.google.cast.receiver",
		payload:     b,
	}
}

type getMediaStatusRequest struct {
	*requestWithTypeAndID
	MediaSessionID int `json:"mediaSessionId"`
}

func newGetMediaStatusRequest(s, d string, mediaSessionID, id int) *pkg {
	b, _ := json.Marshal(&getMediaStatusRequest{
		requestWithTypeAndID: &requestWithTypeAndID{
			Type:      "GET_STATUS",
			RequestID: id,
		},
		MediaSessionID: mediaSessionID,
	})
	return &pkg{
		source:      s,
		destination: d,
		namespace:   "urn:x-cast:com.google.cast.media", // "urn:x-cast:com.google.cast.receiver",
		payload:     b,
	}
}

type launchRequest struct {
	Type      string `json:"type"`
	RequestID int    `json:"requestId"`
	AppID     string `json:"appId"`
}

func newLaunchRequest(s, d string, id int, appID string) *pkg {
	b, _ := json.Marshal(&launchRequest{Type: "LAUNCH", RequestID: id, AppID: appID})
	return &pkg{
		source:      s,
		destination: d,
		namespace:   "urn:x-cast:com.google.cast.receiver",
		payload:     b,
	}
}

type track struct {
	Name             string   `json:"name"`
	TrackID          int      `json:"trackId"`
	Type             string   `json:"type"`
	TrackContentID   string   `json:"trackContentId"`
	TrackContentType string   `json:"trackContentType"`
	Subtype          string   `json:"subtype"`
	Roles            []string `json:"roles"`
	Language         string   `json:"language"`
	IsInband         bool     `json:"isInband"`
}

func newSubtitleTrack(subtitleURI string) track {
	return track{
		Name:             "Subtitle Track", // TODO: configurable
		TrackID:          0,                // TODO: configurable
		Type:             "TEXT",           // TODO: configurable
		TrackContentID:   subtitleURI,
		TrackContentType: "text/vtt",           // TODO: configurable
		Subtype:          "SUBTITLES",          // TODO: configurable
		Roles:            []string{"subtitle"}, // TODO: configurable
		Language:         "en",                 // TODO: configurable
		IsInband:         false,                // TODO: configurable
	}
}

type textTrackStyle struct {
	BackgroundColor   string `json:"backgroundColor"`
	EdgeColor         string `json:"edgeColor"`
	EdgeType          string `json:"edgeType"`
	FontFamily        string `json:"fontFamily"`
	FontGenericFamily string `json:"fontGenericFamily"`
	FontScale         int    `json:"fontScale"`
	FontStyle         string `json:"fontStyle"`
	ForegroundColor   string `json:"foregroundColor"`
}

var defaultTextTrackStyle = textTrackStyle{
	BackgroundColor:   "#00000000",
	EdgeColor:         "#000000FF",
	EdgeType:          "OUTLINE",
	FontFamily:        "ARIAL",
	FontGenericFamily: "SANS_SERIF",
	FontScale:         1,
	FontStyle:         "NORMAL",
	ForegroundColor:   "#FFFFFFFF",
}

type mediaRequest struct {
	ContentID      string         `json:"contentId"`
	StreamType     string         `json:"streamType"`
	ContentType    string         `json:"contentType"`
	Tracks         []track        `json:"tracks"`
	TextTrackStyle textTrackStyle `json:"textTrackStyle"`
}

type loadRequest struct {
	Type           string       `json:"type"`
	RequestID      int          `json:"requestId"`
	ActiveTrackIDs []int        `json:"activeTrackIds"`
	Media          mediaRequest `json:"media"`
}

func newLoadRequest(s, d, uriVideo, uriSubtitle string, id int) *pkg {
	var tracks []track
	var trackStyle textTrackStyle
	var activeTrackIds []int
	if uriSubtitle != "" {
		tracks = []track{newSubtitleTrack(uriSubtitle)}
		trackStyle = defaultTextTrackStyle // TODO: configurable
		activeTrackIds = []int{0}          // TODO: configurable
	}
	b, _ := json.Marshal(&loadRequest{
		Type:           "LOAD",
		RequestID:      id,
		ActiveTrackIDs: activeTrackIds,
		Media: mediaRequest{
			ContentID:      uriVideo,
			StreamType:     "BUFFERED",  // TODO: configurable
			ContentType:    "video/mp4", // TODO: configurable
			Tracks:         tracks,
			TextTrackStyle: trackStyle,
		},
	})
	return &pkg{
		source:      s,
		destination: d,
		namespace:   "urn:x-cast:com.google.cast.media",
		payload:     b,
	}
}

type mediaSessionRequest struct {
	Type           string `json:"type"`
	RequestID      int    `json:"requestId"`
	MediaSessionID int    `json:"mediaSessionId"`
}

func newMediaSessionRequest(s, d, t string, id, mediaSessionID int) *pkg {
	b, _ := json.Marshal(&mediaSessionRequest{Type: t, RequestID: id, MediaSessionID: mediaSessionID})
	return &pkg{
		source:      s,
		destination: d,
		namespace:   "urn:x-cast:com.google.cast.media",
		payload:     b,
	}
}

type seekMediaSessionRequest struct {
	*mediaSessionRequest
	CurrentTime float64 `json:"currentTime"`
}

func newSeekMediaSessionRequest(s, d string, id, mediaSessionID int, currentTime float64) *pkg {
	b, _ := json.Marshal(&seekMediaSessionRequest{
		mediaSessionRequest: &mediaSessionRequest{
			Type:           "SEEK",
			RequestID:      id,
			MediaSessionID: mediaSessionID,
		},
		CurrentTime: currentTime,
	})
	return &pkg{
		source:      s,
		destination: d,
		namespace:   "urn:x-cast:com.google.cast.media",
		payload:     b,
	}
}

type receiverStatus struct {
	RequestID int    `json:"requestId"`
	Type      string `json:"type"`
	Status    struct {
		Applications []struct {
			AppID             string `json:"appId"`
			AppType           string `json:"appType"`
			DisplayName       string `json:"displayName"`
			IconURL           string `json:"iconUrl"`
			IsIdleScreen      bool   `json:"isIdleScreen"`
			LaunchedFromCloud bool   `json:"launchedFromCloud"`
			Namespaces        []struct {
				Name string `json:"name"`
			} `json:"namespaces"`
			SessionID      string `json:"sessionId"`
			StatusText     string `json:"statusText"`
			TransportID    string `json:"transportId"`
			UniversalAppID string `json:"universalAppId"`
		} `json:"applications"`
		UserEq struct{} `json:"userEq"`
		Volume struct {
			ControlType  string  `json:"controlType"`
			Level        float64 `json:"level"`
			Muted        bool    `json:"muted"`
			StepInterval float64 `json:"stepInterval"`
		} `json:"volume"`
	} `json:"status"`
}

type media struct {
	ContentID   string `json:"contentId"`
	StreamType  string `json:"streamType"`
	ContentType string `json:"contentType"`
	Tracks      []struct {
		Name             string   `json:"name"`
		TrackID          int      `json:"trackId"`
		Type             string   `json:"type"`
		TrackContentID   string   `json:"trackContentId"`
		TrackContentType string   `json:"trackContentType"`
		Subtype          string   `json:"subtype"`
		Roles            []string `json:"roles"`
		Language         string   `json:"language"`
		IsInband         bool     `json:"isInband"`
	} `json:"tracks"`
	TextTrackStyle struct {
		BackgroundColor   string `json:"backgroundColor"`
		EdgeColor         string `json:"edgeColor"`
		EdgeType          string `json:"edgeType"`
		FontFamily        string `json:"fontFamily"`
		FontGenericFamily string `json:"fontGenericFamily"`
		FontScale         int    `json:"fontScale"`
		FontStyle         string `json:"fontStyle"`
		ForegroundColor   string `json:"foregroundColor"`
	} `json:"textTrackStyle"`
	MediaCategory string  `json:"mediaCategory"`
	Duration      float64 `json:"duration"`
}

type mediaStatus struct {
	RequestID int    `json:"requestId"`
	Type      string `json:"type"`
	Status    []struct {
		MediaSessionID         int     `json:"mediaSessionId"`
		PlaybackRate           int     `json:"playbackRate"`
		PlayerState            string  `json:"playerState"`
		CurrentTime            float64 `json:"currentTime"`
		SupportedMediaCommands int     `json:"supportedMediaCommands"`
		Volume                 struct {
			Level float64 `json:"level"`
			Muted bool    `json:"muted"`
		} `json:"volume"`
		ActiveTrackIDs []int `json:"activeTrackIds"`
		Media          media `json:"media"`
		CurrentItemID  int   `json:"currentItemId"`
		IdleReason     string
		ExtendedStatus struct {
			PlayerState    string `json:"playerState"`
			Media          media  `json:"media"`
			MediaSessionId int    `json:"mediaSessionId"`
		} `json:"extendedStatus"`
		Items []struct {
			ItemID         int   `json:"itemId"`
			Media          media `json:"media"`
			ActiveTrackIDs []int `json:"activeTrackIds"`
			OrderID        int   `json:"orderId"`
		} `json:"items"`
		RepeatMode string `json:"repeatMode"`
	} `json:"status"`
}
