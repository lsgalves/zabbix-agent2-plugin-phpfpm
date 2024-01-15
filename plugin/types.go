package plugin

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"git.zabbix.com/ap/plugin-support/zbxerr"
)

type unixTime int64
type listenType int

type ErrorMessage struct {
	Message string `json:"message"`
}

const (
	Socket listenType = iota
	Http
)

type Pool struct {
	Name       string
	Listen     string
	ListenType listenType
	ConfigPath string
}

type PoolDiscovery struct {
	Name       string     `json:"{#NAME}"`
	Listen     string     `json:"{#LISTEN}"`
	ListenType listenType `json:"{#LISTEN_TYPE}"`
	ConfigPath string     `json:"{#CONFIG_PATH}"`
}

type PoolStatus struct {
	Pool               string   `json:"pool"`
	ProcessManager     string   `json:"process manager"`
	StartTime          unixTime `json:"start time"`
	StartSince         int      `json:"start since"`
	AcceptedConn       int      `json:"accepted conn"`
	ListenQueue        int      `json:"listen queue"`
	MaxListenQueue     int      `json:"max listen queue"`
	ListenQueueLen     int      `json:"listen queue len"`
	IdleProcesses      int      `json:"idle processes"`
	ActiveProcesses    int      `json:"active processes"`
	TotalProcesses     int      `json:"total processes"`
	MaxActiveProcesses int      `json:"max active processes"`
	MaxChildrenReached int      `json:"max children reached"`
	SlowRequests       int      `json:"slow requests"`
}

// Function to get status of a pool
func (p *Pool) GetStatus(timeout int) (*PoolStatus, error) {
	var (
		resp *http.Response
		err  error
	)

	switch p.ListenType {
	case Socket:
		resp, err = fcgiRequest(p.Listen, "/status", timeout)
	case Http:
		resp, err = httpRequest(p.Listen, "/status", timeout)
	}

	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	var ps PoolStatus

	if err = json.Unmarshal(body, &ps); err != nil {
		return nil, zbxerr.ErrorCannotUnmarshalJSON.Wrap(err)
	}

	return &ps, nil
}

// Function to ping a pool
func (p *Pool) Ping(timeout int) (int, error) {
	var (
		resp *http.Response
		err  error
	)

	switch p.ListenType {
	case Socket:
		resp, err = fcgiRequest(p.Listen, "/ping", timeout)
	case Http:
		resp, err = httpRequest(p.Listen, "/ping", timeout)
	}

	if err != nil {
		return 0, zbxerr.ErrorCannotFetchData.Wrap(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	if strings.TrimSpace(string(body)) == "pong" {
		return 1, nil
	} else {
		return 0, nil
	}
}
