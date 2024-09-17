# WCL Test Application

This project consists of a backend server, a frontend application, and an admin interface.

## Prerequisites
- Go
  - mingw (for Windows)
    - https://scoop.sh/
    - scoop install mingw
- Bun (or Node.js)
- Rust

## Additional dependencies
- linux:
  - pkg-config
  - webkit2gtk-4.1
- windows:
  - pkg-config
  - mingw
  - git bash

## Packaging the application
- windows (application):
  - the end user's system must have WebView2.dll (better to ship it with the application. In the same directory as the application)
  - the end user's system must have urita.dll in the same directory as the application
  - an additional .env file can be placed in the application directory to override variables
    - SERVER_URL: the uri of the server
- linux (server):
  - must ship a .env with the following variables:
    - MONGODB_URI: the uri of the mongodb server
    - DB_NAME: the name of the database
    - CORS_ALLOW_ORIGINS: the origins that are allowed to access the server

## Setup
```bash
./run.sh setup
```

## Build urita
```bash
./run.sh build-urita
```

## Run Vite for frontend
```bash
./run.sh web-dev
```

## Run backend
```bash
./run.sh server
```

## Run application
```bash
./run.sh app
```

## Building
```bash
# for windows (git bash)
./run.sh build-windows-app
./run.sh build-windows-server

# for linux
./run.sh build-app
./run.sh build-server
```

## Cross-compiling for Windows from Linux
first you need a working cross compiler setup. one example of how to do this is in flake.nix.

```bash
nix develop .#windows -c run build-windows-app
```
