.PHONY: all build clean build-all install test

# 默认构建当前平台
all: build

# 构建当前平台
build:
	go build -o netspeed ./cmd/netspeed

# 清理编译产物
clean:
	rm -f netspeed netspeed.exe
	rm -rf bin/

# 跨平台编译所有目标
build-all:
	mkdir -p bin
	@echo "构建 Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build -o bin/netspeed-linux-amd64 ./cmd/netspeed
	@echo "构建 Linux ARM64..."
	GOOS=linux GOARCH=arm64 go build -o bin/netspeed-linux-arm64 ./cmd/netspeed
	@echo "构建 macOS AMD64..."
	GOOS=darwin GOARCH=amd64 go build -o bin/netspeed-darwin-amd64 ./cmd/netspeed
	@echo "构建 macOS ARM64 (M1/M2)..."
	GOOS=darwin GOARCH=arm64 go build -o bin/netspeed-darwin-arm64 ./cmd/netspeed
	@echo "构建 Windows AMD64..."
	GOOS=windows GOARCH=amd64 go build -o bin/netspeed-windows-amd64.exe ./cmd/netspeed
	@echo "构建完成！二进制文件位于 bin/ 目录"

# 安装到系统（仅 Linux/macOS）
install: build
	@echo "安装到 /usr/local/bin/..."
	sudo cp netspeed /usr/local/bin/
	@echo "安装完成！现在可以在任何位置使用 'netspeed' 命令"

# 运行测试
test:
	go test -v ./...

# 运行示例
example-test:
	./netspeed -test

example-ip:
	./netspeed -ip

example-watch:
	./netspeed -test -watch 30

example-config:
	./netspeed -test -config sites.example.json
