package plugin

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"git.zabbix.com/ap/plugin-support/zbxerr"
	fcgiclient "github.com/tomasen/fcgi_client"
)

func fcgiRequest(socket string, path string, timeout int) (*http.Response, error) {
	opts := make(map[string]string)
	opts["SCRIPT_NAME"] = path
	opts["SCRIPT_FILENAME"] = path
	opts["REQUEST_METHOD"] = "GET"
	opts["QUERY_STRING"] = "json"

	fcgi, err := fcgiclient.DialTimeout("unix", socket, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	resp, err := fcgi.Get(opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func httpRequest(url string, path string, timeout int) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url + path)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getPools() ([]Pool, error) {
	var pools []Pool

	cmd := exec.Command("ps", "ax")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
	}

	reFPMConf := regexp.MustCompile(`\(([^)]+)\)`)
	rePoolName := regexp.MustCompile(`^\s*\[([^\]]+)\]\s*$`)
	reListen := regexp.MustCompile(`^\s*listen\s*=\s*(.*)\s*$`)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "php-fpm: master process") {
			matchPhp := reFPMConf.FindStringSubmatch(line)
			if len(matchPhp) < 2 {
				continue
			}

			phpfpmConfig := matchPhp[1]
			poolsDir := strings.TrimSuffix(phpfpmConfig, filepath.Ext(phpfpmConfig)) + ".d/"

			files, err := os.ReadDir(poolsDir)
			if err != nil {
				return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
			}

			for _, file := range files {
				if file.Type().IsRegular() && filepath.Ext(file.Name()) == ".conf" {
					var pool Pool
					poolConfPath := filepath.Join(poolsDir, file.Name())

					content, err := os.ReadFile(poolConfPath)
					if err != nil {
						return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
					}

					pool.ConfigPath = poolConfPath
					lines := strings.Split(string(content), "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, ";") {
							continue
						}

						matchPool := rePoolName.FindStringSubmatch(line)
						if len(matchPool) == 2 {
							pool.Name = matchPool[1]
						}

						matchListen := reListen.FindStringSubmatch(line)
						if len(matchListen) == 2 {
							pool.Listen = strings.TrimSpace(matchListen[1])
							if strings.Contains(pool.Listen, ".sock") {
								pool.ListenType = Socket
							} else {
								pool.ListenType = Http
							}
						}
					}

					if pool.Name != "" && pool.Listen != "" {
						pools = append(pools, pool)
					}
				}
			}

		}
	}

	return pools, nil
}
