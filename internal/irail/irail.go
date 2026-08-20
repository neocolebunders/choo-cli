// Package irail is a minimal client for the iRail connections API.
package irail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Stop struct {
	Station  string
	Time     time.Time
	Delay    time.Duration
	Platform string
	Canceled bool
}

type Connection struct {
	Departure Stop
	Arrival   Stop
	Transfers int
}

// iRail's wire format: numbers and booleans as strings, times as unix seconds.
type wireStop struct {
	Station  string `json:"station"`
	Time     string `json:"time"`
	Delay    string `json:"delay"`
	Platform string `json:"platform"`
	Canceled string `json:"canceled"`
}

func (w wireStop) stop() Stop {
	sec, _ := strconv.ParseInt(w.Time, 10, 64)
	delay, _ := strconv.Atoi(w.Delay)
	return Stop{
		Station:  w.Station,
		Time:     time.Unix(sec, 0),
		Delay:    time.Duration(delay) * time.Second,
		Platform: w.Platform,
		Canceled: w.Canceled == "1",
	}
}

var errMessages = map[int]string{
	http.StatusNotFound:            "no route found — check the station names",
	http.StatusBadRequest:          "iRail rejected the request — check the station names",
	http.StatusTooManyRequests:     "too many requests — wait a moment and try again",
	http.StatusInternalServerError: "iRail is having trouble, try again later",
	http.StatusServiceUnavailable:  "iRail is down for the moment, try again later",
}

// Connections returns upcoming journeys between two stations.
func Connections(from, to string) ([]Connection, error) {
	q := url.Values{"format": {"json"}, "lang": {"en"}, "from": {from}, "to": {to}}
	resp, err := http.Get("https://api.irail.be/connections/?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var e struct {
			Message string `json:"message"`
		}
		json.NewDecoder(resp.Body).Decode(&e)
		msg, ok := errMessages[resp.StatusCode]
		if !ok {
			msg = "iRail is having trouble (" + resp.Status + ")"
		}
		if e.Message != "" {
			msg += " — " + e.Message
		}
		return nil, fmt.Errorf("%s", msg)
	}
	var r struct {
		Connection []struct {
			Departure wireStop `json:"departure"`
			Arrival   wireStop `json:"arrival"`
			Vias      *struct {
				Number string `json:"number"`
			} `json:"vias"`
		} `json:"connection"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("irail: bad response: %w", err)
	}
	conns := make([]Connection, len(r.Connection))
	for i, c := range r.Connection {
		transfers := 0
		if c.Vias != nil {
			transfers, _ = strconv.Atoi(c.Vias.Number)
		}
		conns[i] = Connection{Departure: c.Departure.stop(), Arrival: c.Arrival.stop(), Transfers: transfers}
	}
	return conns, nil
}
