import { readFile } from 'fs/promises'
import { getArgsFromMessageText, makeItAQuote } from './utils'

const main = async (): Promise<void> => {
  const avatarBuffer = await readFile('./src/assets/default_profile.png')
  const maskBuffer = await readFile('./src/assets/gradient-mask.png')
  // 放一个空的参数
  const args = getArgsFromMessageText('')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试', args)).writeAsync('pngTest/test1.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试', args)).writeAsync('pngTest/test2.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试', args)).writeAsync('pngTest/test3.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试', args)).writeAsync('pngTest/test4.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试', args)).writeAsync('pngTest/test5.png')
}

// eslint-disable-next-line @typescript-eslint/no-floating-promises
main()
