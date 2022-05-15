import { InputFile } from 'grammy'
import Jimp from 'jimp'
// import path from 'path'
import { UltimateTextToImage, VerticalImage } from 'ultimate-text-to-image'

/**
 * 把任意头像、id、文字转成一张图片
 * @param avatar 头像，buffer 类型
 * @param id 会在前面加上一个 @
 * @param text 图片正文
 */
export const makeItAQuote = async (avatarIn: Buffer | string, maskIn: Buffer, idIn: string, textIn: string): Promise<Jimp> => {
  let avatar: Jimp
  // 这个类型判断是为了通过类型检查
  if (typeof avatarIn === 'string') {
    avatar = await Jimp.read(avatarIn)
  } else {
    avatar = await Jimp.read(avatarIn)
  }
  const text = await Jimp.read(await genTextWithIdPic(idIn, textIn))
  const mask = await Jimp.read(maskIn)
  // 把头像大小缩放到文本的高度
  avatar.resize(text.getHeight(), text.getHeight())
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
export const genTextWithIdPic = async (id: string, text: string): Promise<Buffer> => {
  const image = new VerticalImage([
    new UltimateTextToImage(`"${text}"`, {
      maxWidth: 500,
      fontSize: 32,
      fontColor: 0xffffff,
      lineHeight: 36
    }), // 正文部分
    new UltimateTextToImage('@' + id, {
      fontSize: 20,
      fontColor: 0xffffff99,
      align: 'right',
      width: 500,
      marginTop: 32
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
