import axios from 'axios'
import { getLogger } from 'log4js'

export let token: string
export let notifyChatId: string

/**
 * 从环境变量或参数中获取参数
 * @param envVarName 环境变量名
 * @param argPrefix 参数前缀
 * @param nameForLog 在 log 中打印的名称
 * @returns 指定的参数
 */
export const getEnvVarOrArg = async (envVarName: string, argPrefix: string, nameForLog: string): Promise<string> => {
  const logger = getLogger()
  // 尝试从环境变量中获取
  let arg = process.env[envVarName]
  // 尝试从参数中获取
  if (arg === undefined || arg === '') {
    logger.debug(`没有从环境变量中获取到 ${nameForLog}，尝试从参数中获取...`)
    let foundArg = ''
    for (const arg of process.argv) {
      if (arg.startsWith(argPrefix)) {
        foundArg = arg
        break
      }
    }
    if (foundArg !== '') {
      arg = foundArg.substring(argPrefix.length)
    }
  }
  if (arg === undefined || arg === '') {
    logger.debug(`参数中也没有 ${nameForLog}，报错退出中...`)
    throw new Error(`获取 ${nameForLog} 失败，请确保添加了名为 ${envVarName} 的环境变量，或指定 ${argPrefix} 开头的参数`)
  }
  logger.debug(`获取到了 ${nameForLog}`)
  return arg
}

/**
 * 获取 Bot 接口令牌
 * @returns Bot 接口令牌
 */
export const getBotToken = async (): Promise<string> => {
  const logger = getLogger()
  const token = await getEnvVarOrArg('BOT_TOKEN', '--token=', 'Bot 接口令牌')
  // 尝试验证获取到的 token 是有效的
  try {
    const res = await axios.get<{ ok: boolean }>(`https://api.telegram.org/bot${token}/getMe`)
    if (!res.data.ok) {
      logger.debug('token 是无效的')
      throw new Error('Bot 接口令牌无效，请检查 BOT_TOKEN 环境变量或 --token= 开头的参数')
    }
  } catch (e) {
    logger.debug('token 是无效的')
    throw new Error('Bot 接口令牌无效，请检查 BOT_TOKEN 环境变量或 --token= 开头的参数')
  }
  logger.debug('token 是有效的')
  return token
}

// 初始化配置
export const initConfig = async (): Promise<void> => {
  token = await getBotToken()
  notifyChatId = await getEnvVarOrArg('NOTIFY_CHAT_ID', '--notify-chat-id=', '通知聊天 ID')
}
