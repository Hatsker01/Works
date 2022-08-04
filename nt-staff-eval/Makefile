.PHONY: build

CURRENT_DIR=$(shell pwd)
APP=catalog
APP_CMD_DIR=./cmd

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto:
	./script/gen-proto.sh	${CURRENT_DIR}

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...

mig-up:
	./script/migrate-up.sh		${CURRENT_DIR}

mig-down:
	./script/migrate-down.sh		${CURRENT_DIR}

mig-fix:
	./script/migrate-fix.sh		${CURRENT_DIR}

swag-gen:
	echo ${REGISTRY}
	swag init -g api/router.go -o api/docs

run-swag:
	sudo apt-get install -y golang-github-go-openapi-swag-dev
	go get -v -u github.com/swaggo/swag/cmd/swag

run-docker:
	sudo chmod 666 /var/run/docker.sock

build-image:
	docker build -t gitlab.abutech.uz:5050/hr.najottalim.uz/nt-staff-eval .

push-image:
	docker push gitlab.abutech.uz:5050/hr.najottalim.uz/nt-staff-eval