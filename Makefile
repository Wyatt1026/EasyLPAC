APP_NAME := EasyLPAC
APP_CMD := ./cmd/easylpac

.PHONY: run build test fmt package-macos-arm package-macos-x86 clean-macos

run:
	go run $(APP_CMD)

build:
	go build -o $(APP_NAME) $(APP_CMD)

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

package-macos-arm:
	./build/macos/create_dmg_arm.sh

package-macos-x86:
	./build/macos/create_dmg_x86.sh

clean-macos:
	./build/macos/clean.sh
