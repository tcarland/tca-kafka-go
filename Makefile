# Makefile for wrapping golang builds
#

all: gomod tca-kafka-go test

gomod:
	( go mod tidy )

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

