# Makefile for wrapping golang builds
#
RM=rm -rf

all: tca-kafka-go test

tca-kafka-go:
	( cd kafka && go build )
	( cd utils && go build )

test:
	( go test ./utils/ -v )

distclean: clean
clean: 
	@echo

install: 
	@echo

