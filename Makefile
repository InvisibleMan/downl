APP?=downl
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d %H:%M:%S')

clean:
	rm -f ${APP}

build: clean
	go build \
        -ldflags "-s -w \
        -X main.Commit=${COMMIT} \
		-X 'main.BuildTime=${BUILD_TIME}'" \
        -o ${APP}

run: build
	./${APP}

server:
	docker-compose up

down:
	./${APP} 500 http://localhost:8083/fast/file.dat http://localhost:8083/slow/file.dat


test:
	go test -v -race ./...
