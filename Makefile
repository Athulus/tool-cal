NAME = tool-cal

all: install test build

build:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${NAME}
test:
		go test -v
install:
		go mod download
docker:
		docker build -t ${NAME} .

compose:
		docker-compose up

docker-build: build docker		

deploy: docker-build compose