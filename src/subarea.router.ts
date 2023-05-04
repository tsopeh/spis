import { CheerioCrawlingContext, Dataset, RouterHandler } from 'crawlee'
import { createHash } from 'crypto'

export interface SubareaInfo {
  url: string
  title: string
  hash: string
  content: string
}

export const registerSubareaRoute = async (routerRef: RouterHandler<CheerioCrawlingContext>, datasetRef: Dataset<SubareaInfo>) => {

  const processed = await datasetRef.getData()

  routerRef.addHandler('subarea', async ({ $, request }) => {

    const url = request.url
    const title = request.userData.title as string
    const rawText = $('body').prop('innerText')
    if (rawText == null) {
      throw new Error(`Could not extract raw text for "${title}" at ${url}.`)
    }

    const dataToHash = JSON.stringify({ url, title, rawText })
    const hash = createHash('md5').update(dataToHash, 'utf-8').digest('hex')

    const found = processed.items.find(p => p.hash == hash)
    if (found != null) {
      // console.log('hit')
      return
    }

    console.log('missed', title)

    await datasetRef.pushData({
      url,
      title,
      hash,
      content: rawText,
    })

  })

}