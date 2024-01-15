package plugin

import (
	"git.zabbix.com/ap/plugin-support/conf"
	"git.zabbix.com/ap/plugin-support/plugin"
)

// Options is a plugin configuration
type Options struct {
	plugin.SystemOptions `conf:"optional,name=System"`

	Timeout int `conf:"optional,range=1:30"`
}

// Configure implements the Configurator interface.
// Initializes configuration structures.
func (p *Plugin) Configure(global *plugin.GlobalOptions, options interface{}) {
	if err := conf.Unmarshal(options, &p.options); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}

	if p.options.Timeout == 0 {
		p.options.Timeout = global.Timeout
	}

	pools, err := getPools()
	if err != nil {
		p.Errf("cannot get pools: %s", err)
	}
	p.pools = pools
}

// Validate implements the Configurator interface.
// Returns an error if validation of a plugin's configuration is failed.
func (p *Plugin) Validate(options interface{}) error {
	var opts Options

	return conf.Unmarshal(options, &opts)
}
