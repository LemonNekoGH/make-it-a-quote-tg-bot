import { getBotToken, getEnvVarOrArg } from './config'
import log4js from 'log4js'
import { Bot, InputFile } from 'grammy'
import { readFile } from 'fs/promises'
import { jimpToInputFile, makeItAQuote } from './utils'
import { registerFont } from 'ultimate-text-to-image'

const main = async (): Promise<void> => {
  // 配置 logger
  log4js.configure({
    appenders: { default: { type: 'console' } },
    categories: { default: { appenders: ['default'], level: 'debug' } }
  })
  // 获取 token
  const token = await getBotToken()
  const notifyChatId = await getEnvVarOrArg('NOTIFY_CHAT_ID', '--notify=', '启动时通知到的对话 id')
  // 获取 logger
  const logger = log4js.getLogger()
  // 读取默认的头像文件
  const defaultAvatar = await readFile('./src/assets/default_profile.png')
  logger.debug('读取到了默认的头像文件')
  // 读取遮罩文件
  const mask = await readFile('./src/assets/gradient-mask.png')
  logger.debug('读取到了遮罩文件')
  // 读取字体文件
  registerFont('./src/assets/Alibaba-PuHuiTi-Regular.ttf', {
    family: 'AliBabaPuHui'
  })
  logger.debug('已注册字体')
  // 构建一个 bot
  const bot = new Bot(token)
  bot.command('quote', async (ctx) => {
    const msg = ctx.message
    if (typeof msg === 'undefined') {
      logger.error('收到了指令，但是获取不到对话 id')
      return
    }
    const { message_id: messageId } = msg
    const { id: chatId } = ctx.chat
    logger.debug(`收到了来自 ${chatId} 的 quote 指令，消息 id: ${messageId}`)
    // 进行一些错误的判断
    const replyMsg = ctx.message?.reply_to_message
    if (typeof replyMsg === 'undefined') {
      logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 没有获取到被回复的消息`)
      await ctx.reply('你并没有回复任何人哦', {
        reply_to_message_id: ctx.message?.message_id
      })
      return
    }
    const sender = replyMsg.from
    if (typeof sender === 'undefined') {
      logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 没有获取到被回复者`)
      await ctx.reply('你回复的这条消息可能来自一个频道，获取不到作者呢', {
        reply_to_message_id: ctx.message?.message_id
      })
      return
    }
    if (typeof replyMsg.text === 'undefined') {
      logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 没有获取到被回复消息的内容`)
      await ctx.reply('你回复的这条消息没有内容呢', {
        reply_to_message_id: ctx.message?.message_id
      })
      return
    }
    // 被回复者 id
    const username = sender.username ?? 'no_name'
    // 被回复的消息内容
    const text = replyMsg.text
    // 被回复者的头像，这里取第一个
    const avatar = (await ctx.api.getUserProfilePhotos(sender.id)).photos
    let quoted: InputFile | undefined
    if (avatar.length === 0) {
      logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 被回复的消息作者是没有头像的`)
      // 如果没有头像，使用默认头像进行图片的合成
      const res = await makeItAQuote(defaultAvatar, mask, username, text)
      quoted = await jimpToInputFile(res)
    } else {
      // 有头像，使用头像组中的第一个进行图片的合成
      const photo = avatar[0][0]
      // 尝试获取此文件
      const file = await ctx.api.getFile(photo.file_id)
      if (typeof file.file_path === 'undefined') {
        logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 被回复的消息作者头像文件获取到了，但是没有路径`)
        return
      }
      logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 被回复的消息作者头像在 https://api.telegram.org/file/bot${token}/${file.file_path}`)
      const res = await makeItAQuote(`https://api.telegram.org/file/bot${token}/${file.file_path}`, mask, username, text)
      quoted = await jimpToInputFile(res)
    }
    await ctx.replyWithPhoto(quoted, {
      reply_to_message_id: messageId
    })
  })
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
