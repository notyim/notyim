package gaia

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/notyim/gaia/dao"
	"github.com/notyim/gaia/scanner/httpscanner"
	"github.com/notyim/gaia/scanner/tcpscanner"
)

type EventType int

const (
	EventTypeCheckInsert EventType = iota
	EventTypeCheckReplace
	EventTypeCheckDelete
)

const (
	EventTypeRunCheck = iota + 1000
	EventTypeCheckHTTPResult
	EventTypeBeat
	EventTypeCheckTCPResult
)

const (
	EventTypePing = iota + 2000
)

type EventCheckInsert struct {
	EventType EventType
	*dao.Check
}

type EventCheckReplace struct {
	EventType EventType
	*dao.Check
}

type EventCheckDelete struct {
	EventType EventType
	*dao.Check
}

type EventCheckHTTPResult struct {
	EventType EventType
	ID        string
	Agent     string
	Region    string
	Result    *httpscanner.CheckResponse
}

func (e *EventCheckHTTPResult) MetricPayload() (map[string]interface{}, error) {
	return map[string]interface{}{
		"time_NameLookup":    e.Result.Timing.NameLookup,
		"time_Connect":       e.Result.Timing.Connect,
		"time_TLSHandshake":  e.Result.Timing.TLSHandshake,
		"time_StartTransfer": e.Result.Timing.StartTransfer,
		"time_Total":         e.Result.Timing.Total,
		"status_code":        e.Result.Status,
	}, nil
}

func (e *EventCheckHTTPResult) CheckID() string {
	return e.ID
}

func (e *EventCheckHTTPResult) CheckType() string {
	return "http"
}

func (e *EventCheckHTTPResult) QueuePayload() ([]byte, error) {
	payload, err := json.Marshal(e.Result)
	if err != nil {
		return nil, fmt.Errorf("Cannot encode json %w", err)
	}

	return payload, nil
}

type EventCheckTCPResult struct {
	EventType EventType `json:"event_type"`
	ID        string
	Agent     string
	Region    string
	Result    *tcpscanner.CheckResponse
}

func (e *EventCheckTCPResult) MetricPayload() (map[string]interface{}, error) {
	return map[string]interface{}{
		"time_Total": e.Result.Timing.Total,
		"error":      e.Result.Error,
		"port_open":  e.Result.PortOpen,
	}, nil
}

func (e *EventCheckTCPResult) QueuePayload() ([]byte, error) {
	payload, err := json.Marshal(e.Result)
	if err != nil {
		return nil, fmt.Errorf("Cannot encode json %w", err)
	}

	return payload, nil
}

func (e *EventCheckTCPResult) CheckID() string {
	return e.ID
}

func (e *EventCheckTCPResult) CheckType() string {
	return "tcp"
}

type EventCheckBeat struct {
	EventType EventType
	ID        string
	Action    string
	BeatAt    time.Time
}

type EventRunCheck struct {
	EventType EventType
	ID        string
}

type EventPing struct {
	EventType EventType
	At        time.Time
}

func NewEventPing() *EventPing {
	return &EventPing{EventType: EventTypePing}
}