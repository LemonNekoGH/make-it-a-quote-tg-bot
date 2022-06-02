import { readFile } from 'fs/promises'
import { getLogger } from 'log4js'
import { registerFont } from 'ultimate-text-to-image'

export let defaultAvatar: Buffer
export let mask: Buffer

export const loadAssets = async (): Promise<void> => {
  // 获取 logger
  const logger = getLogger()
  // 读取默认的头像文件
  defaultAvatar = await readFile('./src/assets/default_profile.png')
  logger.debug('读取到了默认的头像文件')
  // 读取遮罩文件
  mask = await readFile('./src/assets/gradient-mask.png')
  logger.debug('读取到了遮罩文件')
  // 读取字体文件
  registerFont('./src/assets/Alibaba-PuHuiTi-Regular.ttf', {
    family: 'AliBabaPuHui'
  })
  logger.debug('已注册字体文件')
}
