package plugin

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.zabbix.com/ap/plugin-support/plugin"
	"git.zabbix.com/ap/plugin-support/zbxerr"
)

const Name = "PHPFPM"

type Plugin struct {
	plugin.Base
	options Options
	pools   []Pool
}

var Impl Plugin

func (p *Plugin) Export(key string, rawParams []string, _ plugin.ContextProvider) (result interface{}, err error) {
	params, _, _, err := metrics[key].EvalParams(rawParams, nil)
	if err != nil {
		return nil, err
	}

	switch key {
	case keyPools:
		pools, err := getPools()
		if err != nil {
			return nil, err
		}

		res, err := json.Marshal(pools)
		if err != nil {
			return nil, zbxerr.ErrorCannotMarshalJSON.Wrap(err)
		}

		result = string(res)

	case keyPoolStatus:
		pool := p.findPoolByName(params["Pool"])
		if pool == nil {
			notFoundError := errors.New(fmt.Sprintf("Pool %s not found", params["Pool"]))
			return nil, zbxerr.ErrorEmptyResult.Wrap(notFoundError)
		}

		ps, err := pool.GetStatus(p.options.Timeout)
		if err != nil {
			return nil, err
		}

		res, err := json.Marshal(ps)
		if err != nil {
			return nil, zbxerr.ErrorCannotMarshalJSON.Wrap(err)
		}

		result = string(res)

	case keyPoolPing:
		pool := p.findPoolByName(params["Pool"])
		if pool == nil {
			notFoundError := errors.New(fmt.Sprintf("Pool %s not found", params["Pool"]))
			return nil, zbxerr.ErrorEmptyResult.Wrap(notFoundError)
		}

		res, err := pool.Ping(p.options.Timeout)
		if err != nil {
			return nil, err
		}
		result = res

	case keyPoolDiscovery:
		pools, err := getPools()
		if err != nil {
			return nil, err
		}

		var discoveryPools []PoolDiscovery
		for _, pool := range pools {
			discoveryPools = append(discoveryPools, PoolDiscovery{
				Name:       pool.Name,
				Listen:     pool.Listen,
				ListenType: pool.ListenType,
				ConfigPath: pool.ConfigPath,
			})
		}

		res, err := json.Marshal(discoveryPools)
		if err != nil {
			return nil, zbxerr.ErrorCannotMarshalJSON.Wrap(err)
		}

		result = string(res)
	}

	return result, nil
}

func (p *Plugin) findPoolByName(name string) *Pool {
	for _, pool := range p.pools {
		if pool.Name == name {
			return &pool
		}
	}
	return nil
}
