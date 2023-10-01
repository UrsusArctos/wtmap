package wtapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
)

const (
	// WT API endpoints
	apiMapInfo = "map_info.json"
	apiMapObj  = "map_obj.json"
	apiMapFile = "map.img?gen=%d"
	// WT states
	WTClientStateNotStarted = iota
	WTClientStateIdle       = iota
	WTClientStateInSession  = iota
)

type (
	TMapInfo struct {
		GridSize      [2]float64 `json:"grid_size"`
		GridSteps     [2]float64 `json:"grid_steps"`
		GridZero      [2]float64 `json:"grid_zero"`
		MapGeneration int64      `json:"map_generation"`
		MapMax        [2]float64 `json:"map_max"`
		MapMin        [2]float64 `json:"map_min"`
		Valid         bool       `json:"valid"`
	}

	TMapObjects []TMapObject

	TObjectColor = []uint8

	TMapObject struct {
		Type        string       `json:"type"`
		HTMLColor   string       `json:"color"`
		ObjectColor TObjectColor `json:"color[]"`
		Blink       int          `json:"blink"`
		Icon        string       `json:"icon"`
		IconBg      string       `json:"icon_bg"`
		// Coordinates of a point object (all vehicles)
		// Generally in the form of float values between 0 and 1 relative to the current (!) generation map
		// NOTE: Objects outside current generation map (such as airfields, aircraft spawns etc) may have coordinates beyond [0..1]
		// NOTE: New map generation may occur each time another vehicle is selected in the lobby lineup
		X *float64 `json:"x,omitempty"`
		Y *float64 `json:"y,omitempty"`
		// Course unit vector for player, aircraft and aircraft respawn points
		// These appear to be needing normalization first (i.e. their values sometimes reach beyond [-1..1])
		DX *float64 `json:"dx,omitempty"`
		DY *float64 `json:"dy,omitempty"`
		// Starting coordinates of a directional object with length (airfields, carriers)
		SX *float64 `json:"sx,omitempty"`
		SY *float64 `json:"sy,omitempty"`
		// Ending coordinates of the same
		EX *float64 `json:"ex,omitempty"`
		EY *float64 `json:"ey,omitempty"`
	}

	TWTAPIClient struct {
		hostName   string
		httpClient http.Client
		mapInfo    TMapInfo
	}
)

// Color adapter
func (mobj TMapObject) AdaptColor() color.RGBA {
	return color.RGBA{R: mobj.ObjectColor[0], G: mobj.ObjectColor[1], B: mobj.ObjectColor[2], A: 0xFF}
}

// WTAPIC constructor
func NewClient(hName string) TWTAPIClient {
	// Instantiate HTTP client, prepare damage log
	return TWTAPIClient{hostName: hName, httpClient: *http.DefaultClient}
}

// WTAPIC destructor
func (wtapic *TWTAPIClient) Close() {
	wtapic.httpClient.CloseIdleConnections()
}

// WTAPI URL formatter helper
func (wtapic TWTAPIClient) formatAPIURL(apiMethod string) string {
	return fmt.Sprintf("http://%s:8111/%s", wtapic.hostName, apiMethod)
}

// Generalized WTAPI query
func (wtapic TWTAPIClient) queryWTAPI(endPoint string) ([]byte, error) {
	qBody := &bytes.Buffer{}
	qReq, qErr := http.NewRequest("GET", wtapic.formatAPIURL(endPoint), qBody)
	if (qReq != nil) && (qErr == nil) {
		defer qReq.Body.Close()
		qResp, qErr := wtapic.httpClient.Do(qReq)
		if (qResp != nil) && (qErr == nil) {
			defer qResp.Body.Close()
			qRawResp, qErr := io.ReadAll(qResp.Body)
			if qErr == nil {
				return qRawResp, nil
			}
			return nil, qErr
		}
		return nil, qErr
	}
	return nil, qErr
}

// WTAPIC queries
// Query: determine if WT Client is running and if yes, if it is in game session
func (wtapic *TWTAPIClient) GetWTClientState() int {
	qString, err := wtapic.queryWTAPI(apiMapInfo)
	if err == nil {
		qDec := json.NewDecoder(bytes.NewReader(qString))
		qDec.UseNumber()
		decodErr := qDec.Decode(&wtapic.mapInfo)
		if decodErr == nil {
			if wtapic.mapInfo.Valid {
				return WTClientStateInSession
			} else {
				return WTClientStateIdle
			}
		}
	}
	return WTClientStateNotStarted
}

// Query: get current map generation
func (wtapic TWTAPIClient) GetMapGeneration() int64 {
	return wtapic.mapInfo.MapGeneration
}

// Query: get current map
func (wtapic TWTAPIClient) GetMapFile() ([]byte, error) {
	return wtapic.queryWTAPI(fmt.Sprintf(apiMapFile, wtapic.GetMapGeneration()))
}

// Query: retrieve map objects (only works while in session)
func (wtapic TWTAPIClient) GetMapObjects() (mapObj TMapObjects) {
	qString, err := wtapic.queryWTAPI(apiMapObj)
	if err == nil {
		qDec := json.NewDecoder(bytes.NewReader(qString))
		qDec.UseNumber()
		decodErr := qDec.Decode(&mapObj)
		if decodErr == nil {
			return mapObj
		}
	}
	return TMapObjects{}
}
