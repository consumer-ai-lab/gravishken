#!/usr/bin/env bash

PROJECT_ROOT="$(git rev-parse --show-toplevel)"

# BUILD_MODE="PROD"
BUILD_MODE="DEV"

APP_PORT=6200
SERVER_PORT=6201
VARS="-X main.build_mode=$BUILD_MODE -X main.port=$APP_PORT"

web-build() {
  cd $PROJECT_ROOT/frontend

  bun run build
  if [[ -d ../application/dist ]]; then
    rm -rf ../application/dist
  fi
  cp -r ./dist ../application/.
}

# - [webview/webview](https://github.com/webview/webview?tab=readme-ov-file#windows)
#   - NOTE: install WebView2 runtime for < Windows 11
# - [MAYBE: WebView2Loader.dll](https://github.com/webview/webview?tab=readme-ov-file#ms-webview2-loader)
build-windows-app() {
  web-build
  
  cd $PROJECT_ROOT/application
  BUILD_MODE="PROD"
  VARS="-X main.build_mode=$BUILD_MODE -X main.port=$APP_PORT"
  GOOS=windows
  GOARCH=amd64
  CGO_ENABLED=1

  go build -ldflags "$VARS -H windowsgui" -o build/gravtest.exe ./src/.
}

server-run() {
  cd $PROJECT_ROOT/backend
  source ./.env

  go run .
}

web-dev() {
  cd $PROJECT_ROOT/frontend

  bun run dev
}

app() {
  cd $PROJECT_ROOT/application

  go build -ldflags "$VARS" -o build/gravtest ./src/.
  ./build/gravtest $@
}

run() {
  set -e 
  # set -o pipefail

  command="$1"
  if [[ $# > 1 ]]; then
    shift
  fi

  case $command in
    "web-build")
      web-build
    ;;
    "build-windows-app")
      build-windows-app
    ;;
    "server-run")
      server-run
    ;;
    "web-dev")
      web-dev
    ;;
    "app")
      app $@
    ;;
    *)
      echo "unknown command"
    ;;
  esac
}

run $@
