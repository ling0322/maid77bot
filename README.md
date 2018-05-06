# maid77bot

A funny bot in Telegram / 一个有趣的Telegram~~女仆~~机器人

Clone代码

```bash
git clone https://github.com/ling0322/maid77bot.git
cd maid77bot
```

编译

```bash
export GOPATH=$PWD
go get maid77bot/...
go install maid77bot/...
```

运行: 将config_template.yaml复制成config.yaml, 填上Telegram bot的token, 然后即可运行bot: `bin/maid77bot`

```bash
cp config_template.yaml config.yaml
nano config.yaml
bin/maid77bot
```
