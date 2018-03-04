# Borrowed from:
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = icon
BINARY_PATH = bin

# Build the project
all: clean build

build:
	go build ${LDFLAGS} -o ${BINARY_PATH}/${BINARY} . ; \

test:
	mkdir dist
	bin/icon -o dist Identicon

clean:
	-rm -f ${BINARY_PATH}/${BINARY}
	-rm -rf dist

install:
	go install

.PHONY: clean build
