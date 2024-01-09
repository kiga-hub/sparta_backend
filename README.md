# 项目文档

- 后端：用`echo`快速搭建基础restful风格API。
- 网络服务器框架：使用`gnet`网络框架实现数据接收。
- API文档：使用`Swagger`构建自动化文档。
- 配置文件：使用`viper`解析配置文件。
- CLI: 使用`cobra`实现命令行参数。

## 1. 安装说明

- golang版本 >= v1.20

### 2. 源码编译

```bash
# 使用 go.mod

# 安装go依赖包
go list (go mod tidy)

# 编译
./build.sh websocket websocket
```

### 2.2 修改配置文件

```bash
vim websocket.toml
```

### 2.3 启动

```bash
# 运行
./websocket run
```

### 2.4 Docker安装

```bash
# 创建docker镜像
make

# 启动
docker run -it --rm -p 8000:8000 --name websocket websocket:dev
```
