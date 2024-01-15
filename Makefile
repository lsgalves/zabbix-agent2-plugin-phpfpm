.ONESHELL:

PACKAGE=zabbix-agent2-plugin-phpfpm
TOPDIR := $(CURDIR)

build:
	CGO_ENABLED=0 go build -o "$(TOPDIR)/$(PACKAGE)"

clean:
	rm -rf "$(TOPDIR)/$(PACKAGE)"*
	go clean "$(TOPDIR)/..."

format:
	go fmt "$(TOPDIR)/..."

install: build
	mkdir -p /usr/sbin/zabbix-agent2-plugin/
	install -o root -g root -m 755 zabbix-agent2-plugin-phpfpm /usr/sbin/zabbix-agent2-plugin/
	install -o root -g root -m 644 phpfpm.conf /etc/zabbix/zabbix_agent2.d/plugins.d/
	sed -i 's|^# \(Plugins.PHPFPM.System.Path=\)|\1/usr/sbin/zabbix-agent2-plugin/zabbix-agent2-plugin-phpfpm|' /etc/zabbix/zabbix_agent2.d/plugins.d/phpfpm.conf
	systemctl restart zabbix-agent2
