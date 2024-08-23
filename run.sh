#!/usr/bin/env bash

export PROJECT_ROOT="$(git rev-parse --show-toplevel)"

# export BUILD_MODE="PROD"
export BUILD_MODE="DEV"

export APP_PORT=6200
export SERVER_PORT=6201
export DEV_PORT=6202

export SERVER_URL="http://localhost:$SERVER_PORT/"
export VARS="-X main.build_mode=$BUILD_MODE -X main.port=$APP_PORT -X main.server_url=$SERVER_URL"

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
  export BUILD_MODE="PROD"
  # export SERVER_URL=""
  export VARS="-X main.build_mode=$BUILD_MODE -X main.port=$APP_PORT -X main.server_url=$SERVER_URL"
  export GOOS=windows
  export GOARCH=amd64
  export CGO_ENABLED=1

  echo "NOTE: building with SERVER_URL as $SERVER_URL"

  go build -ldflags "$VARS -H windowsgui" -o build/gravtest.exe ./src/.
}

build-server() {
  cd $PROJECT_ROOT/backend
  source ./.env

  export BUILD_MODE="PROD"

  export VARS="-X main.build_mode=$BUILD_MODE"
  go build -ldflags "$VARS" -o build/server ./src/.
}

server() {
  cd $PROJECT_ROOT/backend
  source ./.env

  export VARS="-X main.build_mode=$BUILD_MODE"
  go build -ldflags "$VARS" -o build/server ./src/.
  ./build/server $@
}

web-dev() {
  cd $PROJECT_ROOT/frontend

  bun run dev
}

app() {
  cd $PROJECT_ROOT/application

  mkdir -p ./dist

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
    "build-server")
      build-server
    ;;
    "server")
      server $@
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
