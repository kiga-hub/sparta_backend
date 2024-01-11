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

### third-party library

install cmake: version = 3.16.0

```bash
sudo yum remove cmake
```

```bash
wget https://github.com/Kitware/CMake/releases/download/v3.16.0/cmake-3.16.0.tar.gz
```

decompress

```bash
tar -zxvf cmake-3.16.0.tar.gz
```

make and install

```bash
cd cmake-3.16.0

./bootstrap
make
sudo make install

cmake --version
```

#### Ubuntu
```bash
apt install libopenmpi-dev
```

#### Centos
```bash
wget https://download.open-mpi.org/release/open-mpi/v3.1/openmpi-3.1.0.tar.gz
# or go to src pakcage
tar -zxvf openmpi-3.1.0.tar.gz

# install
cd openmpi-3.1.0/
./configure --prefix=/usr/local/openmpi
make -j 48
make install

# get info
whereis openmpi 
# openmpi: /usr/local/openmpi

# set env
vim ~/.bashrc
#add
export PATH=$PATH:/usr/local/openmpi/bin
export LD_LIBRARY_PATH=/usr/local/openmpi/lib

source ~/.bashrc
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
