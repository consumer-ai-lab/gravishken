#!/usr/bin/env bash

export PROJECT_ROOT="$(git rev-parse --show-toplevel)"

# export BUILD_MODE="PROD"
export BUILD_MODE="DEV"

# command to install webkit on fedora
# sudo dnf install webkit2gtk3-devel

export APP_PORT=6200
export SERVER_PORT=6201
export DEV_PORT=6202

# TODO
# export ADMIN_UI_PORT=6203

export SERVER_URL="http://localhost:$SERVER_PORT"
export VARS="-X main.build_mode=$BUILD_MODE -X main.port=$APP_PORT -X main.server_url=$SERVER_URL"


if command -v bun >/dev/null; then
  runner="bun"
else
  runner="npm"
fi

web-build() {
  cd $PROJECT_ROOT/frontend

  $runner run build
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

  go build -ldflags "$VARS -H windowsgui" -o ../build/gravtest.exe ./src/.
}

build-windows-server() {
  cd $PROJECT_ROOT/backend
  source ./.env

  export BUILD_MODE="PROD"
  export VARS="-X main.build_mode=$BUILD_MODE"
  export GOOS=windows
  export GOARCH=amd64
  export CGO_ENABLED=1

  go build -ldflags "$VARS -H windowsgui" -o ../build/server.exe ./src/.
}

build-server() {
  cd $PROJECT_ROOT/backend
  source ./.env

  export BUILD_MODE="PROD"
  # export SERVER_URL=""

  echo "NOTE: building with SERVER_URL as $SERVER_URL"

  export VARS="-X main.build_mode=$BUILD_MODE"
  go build -ldflags "$VARS" -o ../build/server ./src/.
}

admin-server() {
  cd $PROJECT_ROOT/admin

  $runner run dev
}

server() {
  cd $PROJECT_ROOT/backend
  source ./.env

  export VARS="-X main.build_mode=$BUILD_MODE"
  go build -ldflags "$VARS" -o ../build/server ./src/.
  ../build/server $@
}

web-dev() {
  cd $PROJECT_ROOT/frontend

  $runner run dev
}

app() {
  cd $PROJECT_ROOT/application

  mkdir -p ./dist
  touch ./dist/ignore

  go build -ldflags "$VARS" -o ../build/gravtest ./src/.
  ../build/gravtest $@
}

setup() {
  cd $PROJECT_ROOT/application
  go mod tidy

  cd $PROJECT_ROOT/backend
  go mod tidy

  cd $PROJECT_ROOT/admin
  $runner i

  cd $PROJECT_ROOT/frontend
  $runner i
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
    "build-windows-server")
      build-windows-server
    ;;
    "build-server")
      build-server
    ;;
    "setup")
      setup
    ;;
    "admin-server")
      admin-server $@
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
