from golang:alpine as build

# 定义构建时参数
arg version=0.0.1
# 创建 app 文件夹
run mkdir /app
# 指定 app 文件夹为工作目录
workdir /app
# 将当前文件夹中所有内容复制到镜像中
copy . /app/build
# 执行构建
run cd /app/build && go build -ldflags "-X '$(go list ./... | head -n 1)/internal/config.Env=production' -X '$(go list ./... | head -n 1)/internal/config.Version=$version'" -o bot main.go

# 切换镜像基础
from alpine:latest
# 复制构建产物
copy --from=build /app/build/bot /app/bot

# 设定文件映射
volume ["/etc/bot/config"]
# 设定环境变量
env BOT_CONFIG_PATH="/etc/bot/config/config.yaml"
# 设定入口点
entrypoint ["/app/bot"]