import { readFile } from 'fs/promises'
import { makeItAQuote } from './utils'

const main = async (): Promise<void> => {
  const avatarBuffer = await readFile('./src/assets/default_profile.png')
  const maskBuffer = await readFile('./src/assets/gradient-mask.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试')).writeAsync('pngTest/test1.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试')).writeAsync('pngTest/test2.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试')).writeAsync('pngTest/test3.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试')).writeAsync('pngTest/test4.png')
  await (await makeItAQuote(avatarBuffer, maskBuffer, '测试测试', '测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试，测试测试测试测试测试测试测试测试')).writeAsync('pngTest/test5.png')
}

// eslint-disable-next-line @typescript-eslint/no-floating-promises
main()
