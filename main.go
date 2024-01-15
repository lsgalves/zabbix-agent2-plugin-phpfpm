package main

import (
	"errors"
	"fmt"
	"os"

	"git.zabbix.com/ap/plugin-support/plugin/comms"
	"git.zabbix.com/ap/plugin-support/plugin/container"
	"git.zabbix.com/ap/plugin-support/zbxerr"
	"github.com/lsgalves/zabbix-agent2-plugin-phpfpm/plugin"
)

const COPYRIGHT_MESSAGE = `Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

const (
	PLUGIN_VERSION_MAJOR = 1
	PLUGIN_VERSION_MINOR = 0
	PLUGIN_VERSION_PATCH = 0
	PLUGIN_VERSION_RC    = "beta1"
)

func main() {
	err := comms.HandleFlags(
		plugin.Name,
		os.Args[0],
		COPYRIGHT_MESSAGE,
		PLUGIN_VERSION_RC,
		PLUGIN_VERSION_MAJOR,
		PLUGIN_VERSION_MINOR,
		PLUGIN_VERSION_PATCH,
	)
	if err != nil {
		if !errors.Is(err, zbxerr.ErrorOSExitZero) {
			panic(fmt.Sprintf("failed to handle flags %s", err.Error()))
		}

		return
	}

	h, err := container.NewHandler(plugin.Impl.Name())
	if err != nil {
		panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
	}
	plugin.Impl.Logger = &h

	err = h.Execute()
	if err != nil {
		panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
	}
}
