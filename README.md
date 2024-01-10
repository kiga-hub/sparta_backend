# Project Documentation

- Backend: Quickly build basic restful style API with `echo`.
- Network Server Framework: Implement data reception using the `gnet` network framework.
- API Documentation: Build automated documentation with `Swagger`.
- Configuration File: Parse configuration files using `viper`.
- CLI: Implement command line parameters with `cobra`.

## 1. Installation Instructions

- golang version >= v1.20

### 2. Source Code Compilation

```bash
# Using go.mod

# Install go dependency packages
go list (go mod tidy)

# Compile
./build.sh websocket websocket
```

### 2.2 Modify Configuration File

```bash
vim websocket.toml
```

### 2.3 Start

```bash
# Run
./websocket run
```

### 2.4 Docker Installation

```bash
# Create docker image
make

# Start
docker run -it --rm -p 8000:8000 --name websocket websocket:dev
```

