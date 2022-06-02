import { Bot, Middleware } from 'grammy'
import { getLogger } from 'log4js'

/**
 * 命令参数
 * 从消息中获取
 */
export interface ArgsFromMsg {
  quoteMarkLeft: string // 自定义引号字符，左侧
  quoteMarkRight: string // 自定义引号字符，右侧
  gray: boolean // 是否把头像处理成灰色
}

// 获取 Bot 的 command 方法类型
type GetCommandFn = Bot['command']
export type BotCommandHandler = GetCommandFn extends (f: any, r: infer Rest) => any ? Rest : never
export type CommandCtx = BotCommandHandler extends Middleware<infer Ctx> ? Ctx : never
export type MyHandler = (ctx: CommandCtx) => Promise<void>

/**
 * 命令
 */
export class Commnad {
  name: string // 命令名称
  fn: MyHandler // 命令处理函数

  constructor (name: string, fn: MyHandler) {
    this.name = name
    this.fn = fn
  }

  // 让 Bot 使用此命令
  use = (bot: Bot): void => {
    bot.command(this.name, async (ctx) => {
      const logger = getLogger()
      try {
        const msg = ctx.message
        if (typeof msg === 'undefined') {
          logger.info(`收到 ${this.name} 命令，但获取不到消息内容，对话 ID: ${ctx.chat.id}`)
          return
        }
        logger.info(`收到 ${this.name} 命令，对话 ID: ${ctx.chat.id}，消息 ID: ${msg.message_id}`)
        // 尝试执行命令处理器
        await this.fn(ctx)
      } catch (e) {
        logger.error(`[${this.name} 命令，对话 ID: ${ctx.chat.id}] 处理出错：`, e)
      }
    })
  }
}
