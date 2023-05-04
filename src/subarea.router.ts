import { CheerioCrawlingContext, Dataset, RouterHandler } from 'crawlee'

export const registerSubareaRoute = (routerRef: RouterHandler<CheerioCrawlingContext>, datasetRef: Dataset) => {

  routerRef.addHandler('subarea', async ({ $, request }) => {

    const title = request.userData.title as string
    const rawText = $.text()

    datasetRef.pushData({
      url: request.url,
      title,
      text: `${title}\n${rawText}`,
    })

  })

}