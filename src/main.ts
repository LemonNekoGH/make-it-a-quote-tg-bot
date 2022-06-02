import { initConfig, notifyChatId, token } from './config'
import log4js from 'log4js'
import { Bot } from 'grammy'
import { loadAssets } from './assets'
import { useCommands } from './utils'
import { commands } from './commands'

const main = async (): Promise<void> => {
  // 配置 logger
  log4js.configure({
    appenders: { default: { type: 'console' } },
    categories: { default: { appenders: ['default'], level: 'debug' } }
  })
  // 初始化配置
  await initConfig()
  // 获取 logger
  const logger = log4js.getLogger()
  // 加载资源
  await loadAssets()
  // 构建一个 bot
  const bot = new Bot(token)
  // 使用命令处理器集
  useCommands(bot, commands)
  // 错误处理
  bot.catch(e => {
    const err = e as Error
    logger.error('出现了错误, ' + err.message)
  })
  // eslint-disable-next-line @typescript-eslint/no-floating-promises
  bot.start({
    // 在启动时发送一条通知
    onStart: async (info) => {
      await bot.api.sendMessage(notifyChatId, `make-it-a-quote-bot 已启动，bot 用户名是 @${info.username}`)
    }
  })
  logger.info('QuoteMaker Bot 已启动，按下 Ctrl + C 停止运行')
}

// eslint-disable-next-line @typescript-eslint/no-floating-promises
main()
