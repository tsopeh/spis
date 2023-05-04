import { CheerioCrawler, CheerioCrawlingContext, createCheerioRouter, Dataset, RouterHandler } from 'crawlee'
import { registerMenuRoute } from './menu.router.js'
import { registerSubareaRoute, SubareaInfo } from './subarea.router.js'
import { createRandomUniqueKey } from './utils.js'

export const router: RouterHandler<CheerioCrawlingContext> = createCheerioRouter()

// Menu
registerMenuRoute(router)

// Subareas
const subareaHashesDataset = await Dataset.open<SubareaInfo>('subarea')
await registerSubareaRoute(router, subareaHashesDataset)

const crawler = new CheerioCrawler({
  requestHandler: router,
  maxRequestsPerCrawl: 200,
  // maxConcurrency: 1,
})

await crawler.run([
  {
    label: 'menu',
    url: 'https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu',
    uniqueKey: createRandomUniqueKey('menu', 0),
  },
])
