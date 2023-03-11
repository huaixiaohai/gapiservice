# >>>>>>>>>>> 自定义常量 >>>>>>>>>>>>>
# 定义项目基本信息
COMMONENVVAR      ?= GOOS=linux GOARCH=amd64
BUILDENVVAR       ?= CGO_ENABLED=0
BIN_DIR        ?= $(CURDIR)/bin

# >>>>>>>>>>> 必须包含的命令 >>>>>>>>>

# 构建并编译出静态可执行文件
all: linux_build

build:
	go mod tidy
	go build -o $(BIN_DIR)/main

run:
	go run main.go

# 交叉编译出linux下的静态可执行文件build
linux_build:
	$(COMMONENVVAR) $(BUILDENVVAR) make build