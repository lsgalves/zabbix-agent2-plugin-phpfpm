package plugin

import (
	"git.zabbix.com/ap/plugin-support/metric"
	"git.zabbix.com/ap/plugin-support/plugin"
)

const (
	keyPools         = "phpfpm.pools"
	keyPoolStatus    = "phpfpm.pool_status"
	keyPoolPing      = "phpfpm.pool_ping"
	keyPoolDiscovery = "phpfpm.pool.discovery"
)

var paramPool = metric.NewParam("Pool", "Pool name for which the information is needed.").SetRequired()

var metrics = metric.MetricSet{
	keyPools:         metric.New("Returns a list of pools.", nil, false),
	keyPoolDiscovery: metric.New("Returns a list of pools, used for low-level discovery.", nil, false),
	keyPoolStatus: metric.New("Returns realtime status for a given pool.",
		[]*metric.Param{paramPool}, false),
	keyPoolPing: metric.New("Returns health for a given pool.",
		[]*metric.Param{paramPool}, false),
}

func init() {
	plugin.RegisterMetrics(&Impl, Name, metrics.List()...)
}
