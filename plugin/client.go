package plugin

import (
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"git.zabbix.com/ap/plugin-support/zbxerr"
	fcgiclient "github.com/tomasen/fcgi_client"
	"gopkg.in/ini.v1"
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
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "php-fpm: master process") {
			matchPhp := reFPMConf.FindStringSubmatch(line)
			if len(matchPhp) < 2 {
				continue
			}

			phpfpmConfig := strings.TrimSpace(matchPhp[1])
			phpCfg, err := ini.Load(phpfpmConfig)
			if err != nil {
				return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
			}

			sections := phpCfg.Sections()
			for _, section := range sections {

				if section.Name() == "DEFAULT" || section.Name() == "global" {
					hasInclude := section.HasKey("include")
					if hasInclude {
						matches, err := filepath.Glob(section.Key("include").Value())
						if err != nil {
							return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
						}

						for _, match := range matches {
							includeFilePath := strings.TrimSpace(match)
							poolCfg, err := ini.Load(includeFilePath)
							if err != nil {
								return nil, zbxerr.ErrorCannotFetchData.Wrap(err)
							}

							for _, poolSection := range poolCfg.Sections() {
								if poolSection.Name() == "DEFAULT" || poolSection.Name() == "global" {
									continue
								}

								pool := extractPoolFromConfig(includeFilePath, poolSection.Name(), poolSection.Keys())
								if pool.Name != "" && pool.Listen != "" {
									pools = append(pools, pool)
								}
							}
						}
					}
					continue
				} else {
					pool := extractPoolFromConfig(phpfpmConfig, section.Name(), section.Keys())
					if pool.Name != "" && pool.Listen != "" {
						pools = append(pools, pool)
					}
				}
			}
		}
	}

	return pools, nil
}

func extractPoolFromConfig(path string, name string, keys []*ini.Key) Pool {
	var pool Pool
	pool.ConfigPath = path
	pool.Name = name
	for _, key := range keys {
		switch key.Name() {
		case "prefix":
			prefix := key.Value()
			if strings.Contains(prefix, "$pool") {
				prefix = strings.Replace(prefix, "$pool", pool.Name+"/", 1)
			}
			pool.Listen = prefix + pool.Listen

		case "listen":
			pool.Listen = pool.Listen + key.Value()

		}
	}

	if strings.Contains(pool.Listen, ".sock") {
		pool.ListenType = Socket
	} else {
		pool.ListenType = Http
	}

	return pool
}
