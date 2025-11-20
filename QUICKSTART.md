# NetSpeed 快速开始

## 一、安装

### Windows
```powershell
# 下载并运行
.\netspeed.exe -test
```

### Linux/macOS
```bash
# 编译
go build -o netspeed

# 安装到系统
sudo cp netspeed /usr/local/bin/
netspeed -test
```

## 二、基础使用

### 1. 测试网络质量
```bash
netspeed -test
```

### 2. 查看 IP 信息
```bash
netspeed -ip
```

### 3. 持续监控（每30秒）
```bash
netspeed -test -watch 30
```

## 三、代理使用

### SOCKS5 代理
```bash
netspeed -test -proxy socks5://127.0.0.1:1080
```

### HTTP 代理
```bash
netspeed -test -proxy http://proxy.example.com:8080
```

### 检测 VPN 出口 IP
```bash
# 启动 VPN 后
netspeed -ip
```

## 四、自定义配置

### 创建配置文件 sites.json
```json
[
  {"Name": "Google", "URL": "https://www.google.com"},
  {"Name": "Baidu", "URL": "https://www.baidu.com"}
]
```

### 使用配置
```bash
netspeed -test -config sites.json
```

## 五、常用场景

### VPN 连接检测
```bash
# 1. 检测 VPN 前
netspeed -ip

# 2. 启动 VPN

# 3. 检测 VPN 后
netspeed -ip
netspeed -test
```

### 持续监控代理质量
```bash
netspeed -test -proxy socks5://127.0.0.1:1080 -watch 60
```

### 快速测试特定网站
```bash
# 创建 quick.json
echo '[{"Name":"GitHub","URL":"https://github.com"}]' > quick.json

# 测试
netspeed -test -config quick.json
```

## 六、故障排查

### 所有网站都超时
```bash
# 增加超时时间
netspeed -test -timeout 20
```

### 代理不生效
```bash
# 检查代理 URL 格式
netspeed -test -proxy socks5://127.0.0.1:1080  # ✓ 正确
netspeed -test -proxy 127.0.0.1:1080          # ✗ 错误（缺少协议）
```

## 七、获取帮助

```bash
netspeed -help
```
