import { InputFile } from 'grammy'
import Jimp from 'jimp'
import { getLogger } from 'log4js'
// import path from 'path'
import { UltimateTextToImage, VerticalImage } from 'ultimate-text-to-image'

/**
 * 命令参数
 */
export interface Args {
  quoteMarkLeft: string // 自定义引号字符，左侧
  quoteMarkRight: string // 自定义引号字符，右侧
  gray: boolean // 是否把头像处理成灰色
}

/**
 * 把任意头像、id、文字转成一张图片
 * @param avatar 头像，buffer 类型
 * @param id 会在前面加上一个 @
 * @param text 图片正文
 */
export const makeItAQuote = async (avatarIn: Buffer | string, maskIn: Buffer, idIn: string, textIn: string, commandArgs: Args): Promise<Jimp> => {
  const logger = getLogger()
  let avatar: Jimp
  // 这个类型判断是为了通过类型检查
  if (typeof avatarIn === 'string') {
    avatar = await Jimp.read(avatarIn)
  } else {
    avatar = await Jimp.read(avatarIn)
  }
  const text = await Jimp.read(await genTextWithIdPic(idIn, textIn, commandArgs))
  const mask = await Jimp.read(maskIn)
  // 把头像大小缩放到文本的高度
  avatar.resize(text.getHeight(), text.getHeight())
  // 如果需要把头像转成灰色
  if (commandArgs.gray) {
    logger.debug('头像转换成了灰色')
    avatar.greyscale()
  }
  // 把遮罩缩放到文本的高度
  mask.resize(text.getHeight(), text.getHeight())
  // 生成一张宽度是头像和文本宽度之和，高度是文本高度的黑色图片
  const bg = await Jimp.create(text.getHeight() + text.getWidth(), text.getHeight(), 'black')
  // 把头像叠在左边，把文本叠在右边
  return bg.composite(avatar, 0, 0).composite(text, text.getHeight(), 0).composite(mask, 0, 0)
}

/**
 * 用 id 和文字生成一张纯黑底的图片
 * @param id 会在前面加上一个 @
 * @param text 图片正文
 */
export const genTextWithIdPic = async (id: string, text: string, commandArgs: Args): Promise<Buffer> => {
  const logger = getLogger()
  logger.debug(`使用了引号 ${commandArgs.quoteMarkLeft}${commandArgs.quoteMarkRight}`)
  const image = new VerticalImage([
    new UltimateTextToImage(`${commandArgs.quoteMarkLeft}${text}${commandArgs.quoteMarkRight}`, {
      maxWidth: 500,
      fontSize: 32,
      fontColor: 0xffffff,
      lineHeight: 36,
      fontFamily: 'AliBabaPuHui'
    }), // 正文部分
    new UltimateTextToImage('@' + id, {
      fontSize: 20,
      fontColor: 0xffffff99,
      align: 'right',
      width: 500,
      marginTop: 32,
      fontFamily: 'AliBabaPuHui'
    })
  ], {
    margin: 32
  })
  return image.render().toBuffer()
}

/**
 * 把 Jimp 实例导出成 InputFile
 * @param src Jimp 实例
 */
export const jimpToInputFile = async (src: Jimp): Promise<InputFile> => {
  const buffer = await src.getBufferAsync(src.getMIME())
  return new InputFile(buffer)
}

// 从消息文本中获取参数
export const getArgsFromMessageText = ((): ((text: string) => Args) => {
  // 获取参数的方法，不需要外部可以访问
  const getArgValue = (rawArgs: string[], prefix: string): string => {
    for (const arg0 of rawArgs) {
      if (arg0.startsWith(prefix)) {
        return arg0.substring(prefix.length)
      }
    }
    return ''
  }
  // 获取旗标的方法，不需要外部可以访问
  const getFlag = (rawArgs: string[], flag: string): boolean => {
    for (const arg0 of rawArgs) {
      if (arg0 === flag) {
        return true
      }
    }
    return false
  }
  return (text: string): Args => {
    // 默认参数
    const defaultArgs: Args = {
      quoteMarkLeft: '"',
      quoteMarkRight: '"',
      gray: false
    }
    // 参数集合
    const rawArgs = text.split(' ')
    // 是否需要处理成灰色
    defaultArgs.gray = getFlag(rawArgs, 'gray')
    // 是否需要自定义左引号
    defaultArgs.quoteMarkLeft = getArgValue(rawArgs, 'leftMark=') ?? '"'
    // 是否需要自定义右引号
    defaultArgs.quoteMarkRight = getArgValue(rawArgs, 'rightMark=') ?? '"'
    return defaultArgs
  }
})()
