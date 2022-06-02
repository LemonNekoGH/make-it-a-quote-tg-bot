import { InputFile } from 'grammy'
import { getLogger } from 'log4js'
import { defaultAvatar, mask } from './assets'
import { token } from './config'
import { Commnad, MyHandler } from './types'
import { getArgsFromMessageText, jimpToInputFile, makeItAQuote } from './utils'

const logger = getLogger()

// 处理 quote 命令
const handleQuoteCommand: MyHandler = async (ctx) => {
  const msg = ctx.message!
  const { message_id: messageId } = msg
  const { id: chatId } = ctx.chat
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
  const { message_id: replyId } = await ctx.reply('正在进行处理，请稍等...')
  logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 已成功发送“处理中”提示信息`)
  // 进行参数处理
  const args = getArgsFromMessageText(msg.text)
  let username = ''
  let avatarToGet: number
  // 如果是被转发的消息
  if (typeof replyMsg.forward_from !== 'undefined') {
    logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 是一条被转发的消息`)
    // 获取被转发者的用户名
    username = replyMsg.forward_from.username ?? 'no_name'
    avatarToGet = replyMsg.forward_from.id
  } else if (typeof replyMsg.forward_sender_name !== 'undefined') {
    logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 是一条被转发的消息，但是转发人没有允许它人查看被转发者的信息`)
    await ctx.reply('你回复的这条消息是一个被转发的消息，但是受到转发者的隐私设置限制，获取不到被转发人的信息', {
      reply_to_message_id: ctx.message?.message_id
    })
    return
  } else {
    // 获取被回复者的用户名
    username = sender.username ?? 'no_name'
    avatarToGet = sender.id
  }
  // 被回复的消息内容
  const text = replyMsg.text
  // 被回复者的头像，这里取第一个，如果是被转发的消息，获取被转发者的头像
  const avatar = (await ctx.api.getUserProfilePhotos(avatarToGet)).photos
  let quoted: InputFile | undefined
  if (avatar.length === 0) {
    logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 被回复的消息作者是没有头像的`)
    // 如果没有头像，使用默认头像进行图片的合成
    const res = await makeItAQuote(defaultAvatar, mask, username, text, args)
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
    const res = await makeItAQuote(`https://api.telegram.org/file/bot${token}/${file.file_path}`, mask, username, text, args)
    quoted = await jimpToInputFile(res)
  }
  logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 图片处理完成`)
  await ctx.replyWithPhoto(quoted, {
    reply_to_message_id: messageId
  })
  logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 图片发送完成`)
  await ctx.api.deleteMessage(chatId, replyId)
  logger.debug(`[chat: ${chatId}, command: quote, msg: ${messageId}] 提示信息删除完成`)
}

// 导出命令集
export const commands: Commnad[] = [
  new Commnad('quote', handleQuoteCommand)
]
