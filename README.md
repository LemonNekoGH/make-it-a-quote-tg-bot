# make-it-a-quote-tg
一个把被你回复的消息变成图片的 Bot

效果图:  
![Example](./example.jpg)

## 如何部署？
首先把镜像拉下来
```shell
docker pull ghcr.io/lemonnekogh/make-it-a-quote-tg:latest
```
然后运行
```shell
docker run --name <容器名称> \
  -d -it \
  -e BOT_TOKEN=<你的 Bot 接口令牌> \
  -e NOTIFY_CHAT_ID=<启动时要提醒的对话 id> \
  --restart always \
  lemonnekogh/make-it-a-quote-tg
```
加上 `--restart always` 是为了在容器挂掉之后重新启动，不需要的话，去掉就好了
## TodoList
- [ ] 在收到转换指令后，进行一个友好的回复，避免被误认为卡住了
- [ ] 支持 `gray` 参数，收到这个参数时，头像会被处理成灰色
- [ ] 使用 `semantic-release` 语义化镜像版本