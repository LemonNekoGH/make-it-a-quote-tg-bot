// import Jimp from 'jimp'
// import path from 'path'
import { generate } from 'text-to-image'
import { writeFile } from 'fs/promises'

/**
 * 把任意头像、id、文字转成一张图片
 * @param avatar 头像，buffer 类型
 * @param id 会在前面加上一个破折号
 * @param text 图片正文
 */
// @ts-expect-error
const makeItAQuote = (avatar: ArrayBuffer, id: string, text: string): ArrayBuffer => {
  return Buffer.from('')
}

/**
 * 用 id 和文字生成一张纯黑底的图片
 * @param id 会在前面加上一个破折号
 * @param text 图片正文
 */
// @ts-expect-error
const genTextWithIdPic = async (id: string, text: string): ArrayBuffer => {
  // const pic = new Jimp(1200, 630, 0x000000ff)
  // const font = await Jimp.loadFont(path.resolve(__dirname, './assets/Alibaba-PuHuiTi-Regular.ttf'))
  // pic.print(font, 20, 20, `
  //   ${text}
  //   —— @${id}
  // `)
  // await pic.write('test.png')

  const dataUri = await generate(`
  ${text}
  —— @${id}
`, {
    bgColor: 'black',
    textColor: 'white',
    textAlign: 'right'
  })
  const toWrite = Buffer.from(dataUri.split(',')[1], 'base64') // 转成 buffer 以便输出到图片
  await writeFile('test.png', toWrite)
}

genTextWithIdPic('中文测试', 'test')
