APP?=downl
VERSION?=0.1.0
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d %H:%M:%S')

clean:
	rm -f ${APP}

build: clean
	go build \
        -ldflags "-s -w \
		-X main.Version=${VERSION} \
        -X main.Commit=${COMMIT} \
		-X 'main.BuildTime=${BUILD_TIME}'" \
        -o ${APP}

run: build
	./${APP}

test:
	go test -v -race ./...
