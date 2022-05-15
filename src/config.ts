import axios from 'axios'
import { getLogger } from 'log4js'

export const getBotToken = async (): Promise<string> => {
  const logger = getLogger()
  // 尝试从环境变量中获取 token
  let token = process.env.BOT_TOKEN
  // 尝试从参数中获取 token
  if (token === undefined || token === '') {
    logger.debug('没有从环境变量中获取到 token，尝试从参数中获取...')
    let foundArg = ''
    for (const arg of process.argv) {
      if (arg.startsWith('--token=')) {
        foundArg = arg
        break
      }
    }
    if (foundArg !== '') {
      token = foundArg.substring('--token='.length)
    }
  }
  if (token === undefined || token === '') {
    logger.debug('参数中也没有 token，报错退出中...')
    throw new Error('获取 Bot 接口令牌失败，请确保添加了名为 BOT_TOKEN 的环境变量，或指定 --token= 开头的参数')
  }
  logger.debug('获取到了 token')
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
