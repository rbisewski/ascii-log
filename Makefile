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

#
# TODO: test this on more distros
#
install: build
	@echo installing executable file to /usr/bin/ascii-log
	@sudo cp ascii-log /usr/bin/ascii-log
	@sudo mkdir -p /etc/systemd/system/
	@sudo cp ascii-log.service /etc/systemd/system/ascii-log.service
	@sudo systemctl daemon-reload
	@sudo systemctl enable ascii-log
	@sudo systemctl start ascii-log

#
# TODO: test this on more distros
#
uninstall: clean
	@echo removing executable file from /usr/bin/ascii-log
	@sudo systemctl stop ascii-log
	@sudo systemctl disable ascii-log
	@sudo systemctl daemon-reload
	@sudo rm /usr/bin/ascii-log
	@sudo rm /etc/systemd/system/ascii-log.service
