# Zabbix PHP-FPM plugin
This plugin provides a native solution for monitoring multiples PHP-FPM pools without needing a web server.

## Requirements
- Zabbix Agent 2
- Go >= 1.21 (required only to build from source)

## Installation
*Plugins.PHPFPM.System.Path* variable needs to be set in Zabbix agent 2 configuration file with the path to the
PHP-FPM plugin executable. By default the variable is set in **plugin** configuration file *phpfpm.conf* and then
included in the **agent** configuration file *zabbix_agent2.conf*.

E.g:
You should add the following option to the **plugin** configuration file:

    Plugins.PHPFPM.System.Path=/path/to/executable/phpfpm

Then the config file needs to be included in the main Zabbix agent 2 config file via the *Include* command.

E.g:
You should add the following option to the **plugin** configuration file:

    Include=/path/to/config/phpfpm.conf

## PHP-FPM pools config

Ensure that lines from all active PHP-FPM pools are uncommented:

```conf
pm.status_path = /status
ping.path = /ping
```

And that the Zabbix user has access to the php-fpm socket, for this you can add the zabbix user to `listen.acl_users`.
Finally, reload php-fpm config: `systemctl reload php-fpm`.

## Options
PHP-FPM plugin can be executed on its own with these parameters:
* *-h*, *--help* displays help message
* *-V*, *--version* displays the plugin version and license information

## Configuration
The Zabbix Agent's configuration file is used to configure plugins.

**Plugins.PHPFPM.Timeout** — The maximum time in seconds for waiting when a connection has to be established.  
*Default value:* equals the global Timeout configuration parameter.  
*Limits:* 1-30

## Supported keys
**phpfpm.pool_status[\<Pool\>]** — returns returns status for a given pool.
*Parameters:*
Pool (required) — a pool name.

**phpfpm.pool_ping[\<Pool\>]** — pings to php-fpm pool and returns 0 or 1.
*Parameters:*
Pool (required) — a pool name.
*Returns:*
- "1" if the pool is alive;
- "0" if the pool is broken (is returned if there is any error during the test, or the response is not "pong").

**phpfpm.pools[]** — returns the php-fpm pools for all php-fpm running.

**phpfpm.pool.discovery[]** — returns a list of pools, used in low-level discovery rules.

## Troubleshooting
The plugin uses log output of Zabbix agent 2. You can increase debugging level of Zabbix agent 2 if you need more details about the current situation.
