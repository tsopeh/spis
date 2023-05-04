import { CheerioCrawler, CheerioCrawlingContext, createCheerioRouter, Dataset, RouterHandler } from 'crawlee'
import { registerMenuRoute } from './menu.router.js'
import { registerSubareaRoute } from './subarea.router.js'

export const router: RouterHandler<CheerioCrawlingContext> = createCheerioRouter()

// Menu
registerMenuRoute(router)

// Subareas
const subareaDataset = await Dataset.open('subarea')
registerSubareaRoute(router, subareaDataset)

const crawler = new CheerioCrawler({
  requestHandler: router,
})

await crawler.run([
  { label: 'menu', url: 'https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu' },
])
