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

#install: build
#	@echo installing executable file to ${DESTDIR}${PREFIX}/bin
#
#uninstall: clean
#	@echo removing executable file from ${DESTDIR}${PREFIX}/bin
