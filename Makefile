# Version
VERSION = `date +%y.%m`

# If unable to grab the version, default to N/A
ifndef VERSION
    VERSION = "n/a"
endif

#
# Makefile options
#


# State the "phony" targets
.PHONY: all clean build install uninstall


all: build

build:
	@go build

clean:
	@go clean

install: build
	@echo installing executable file to /usr/bin/ascii-log
	@sudo cp ascii-log /usr/bin/ascii-log
	@echo installing cron file to /etc/cron.d/ascii-log
	@sudo cp ascii-log.cron /etc/cron.d/ascii-log

uninstall: clean
	@echo removing executable file from /usr/bin/ascii-log
	@sudo rm /usr/bin/ascii-log
	@echo removing cron file from /etc/cron.d/ascii-log
	@sudo rm /etc/cron.d/ascii-log
