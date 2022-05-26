package hafas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	dateFmt = "02.01.06"
	timeFmt = "15:04"
)

type stboard struct {
	Connections []Conn  `json:"connections"`
	Station     Station `json:"station"`
}

type Station struct {
	AdditionalInformation AdditionalInfo `json:"additionalInformation"`
	Name                  string         `json:"name"`
}

type AdditionalInfo struct {
	AttrHU    bool   `json:"attr_HU"`
	AttrHUVal string `json:"attr_HU_val"`
	AttrMU    bool   `json:"attr_MU"`
}

type Conn struct {
	Attributes          []Attrs       `json:"attributes"`
	AttributesBfr       []Attrs       `json:"attributes_bfr"`
	Cancelled           bool          `json:"cancelled"`
	Date                string        `json:"date"`
	HasGlobalHIMMessage bool          `json:"hasGlobalHIMMessage"`
	HasHIMMessage       bool          `json:"hasHIMMessage"`
	HimMessages         []interface{} `json:"himMessages"`
	Locations           []Location    `json:"locations"`
	MainLocation        Location      `json:"mainLocation"`
	Name                string        `json:"name"`
	Product             Product       `json:"product"`
	TrainInfo           string        `json:"trainInfo"`
}

type Color struct {
	Bg string `json:"bg"`
	Fg string `json:"fg"`
}

type Attrs struct {
	Code          string `json:"code"`
	HightPriority bool   `json:"hightPriority"`
	Priority      string `json:"priority"`
	Text          string `json:"text"`
}

type Product struct {
	Color         Color  `json:"color"`
	Direction     string `json:"direction"`
	DirectionType string `json:"directionType"`
	Icon          string `json:"icon"`
	Line          string `json:"line"`
	LongName      string `json:"longName"`
	Name          string `json:"name"`
	Type          string `json:"type"`
}

func (p Product) String() string {
	line := p.Line
	if p.LongName == "S-Bahn" {
		line = strings.ReplaceAll(p.Line, " ", "")
	}

	return fmt.Sprintf("%s %s %s %s", p.LongName, line, p.DirectionType, html.UnescapeString(p.Direction))
}

func (p Product) Emoji() string {
	if p.Icon == "icon_bus" {
		return "🚍"
	}

	if p.Icon == "icon_train" {
		return "🚆"
	}

	if p.Icon == "icon_tram" {
		return "🚊"
	}

	if p.Icon == "icon_funicular" {
		return "🚞"
	}

	if p.Icon == "icon_cable_car" {
		return "🚡"
	}

	if p.Icon == "icon_boat" {
		return "🛳"
	}

	panic(p.Icon)
}

type RealTime struct {
	Countdown         string `json:"countdown"`
	Date              string `json:"date"`
	Delay             string `json:"delay"`
	HasRealTime       bool   `json:"hasRealTime"`
	IsDelayed         bool   `json:"isDelayed"`
	IsPlatformChanged bool   `json:"isPlatformChanged"`
	Platform          string `json:"platform"`
	Time              string `json:"time"`
}

type Location struct {
	Countdown string        `json:"countdown"`
	Date      string        `json:"date"`
	Location  LocationPoint `json:"location"`
	Platform  string        `json:"platform"`
	RealTime  RealTime      `json:"realTime"`
	Time      string        `json:"time"`
}

type LocationPoint struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	TupelID string `json:"tupelId"`
	Type    string `json:"type"`
	X       int64  `json:"x"`
	Y       int64  `json:"y"`
}

func (c *Client) StBoard(ctx context.Context, stopName, stopID string, ts time.Time, nResults int) ([]Conn, error) {
	ts = ts.In(c.timeZone)

	params := url.Values{}
	params.Add("dirInput", "")
	params.Add("maxJourneys", fmt.Sprint(nResults))
	if stopName != "" {
		params.Add("input", stopName)
	}
	if stopID != "" {
		params.Add("REQStationS0ID", stopID)
	}
	params.Add("time", ts.Format(timeFmt))
	params.Add("date", ts.Format(dateFmt))
	params.Add("boardType", "dep")
	params.Add("start", "1")
	params.Add("tpl", "stbResult2json")

	url := c.buildURL("/stboard.exe/dny", url.Values{})
	resp, err := http.Post(url, "application/x-www-form-urlencoded; charset=UTF-8", bytes.NewBuffer([]byte(params.Encode())))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("non-200 status code")
	}

	ret := stboard{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret.Connections, nil
}
