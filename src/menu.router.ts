import { RequestOptions } from '@crawlee/core/request'
import { CheerioCrawlingContext, RouterHandler } from 'crawlee'
import { createRandomUniqueKey } from './utils.js'

interface MenuResponse {
  id: number
  name: string
  orderBy: number
  level: number
  count: number
  children?: Array<MenuResponse>
}

interface SubareaId {
  id: number
  name: string
}

const extractSubareaIds = (entries: Array<MenuResponse>, acc: Array<SubareaId>): void => {
  for (const entry of entries) {
    if (entry.children == null || entry.children.length == 0) {
      acc.push({ id: entry.id, name: entry.name })
    } else {
      extractSubareaIds(entry.children, acc)
    }
  }
}

export const registerMenuRoute = (routerRef: RouterHandler<CheerioCrawlingContext>) => {

  routerRef.addHandler('menu', async ({ body, crawler, log }) => {
    if (!(body instanceof Buffer)) {
      throw new Error('Expected for `body` to be a Buffer.')
    }

    const menuResponse: Array<MenuResponse> = JSON.parse(body.toString())

    const subareaIdsAcc: Array<SubareaId> = []
    extractSubareaIds(menuResponse, subareaIdsAcc)

    const requests = subareaIdsAcc.slice(0, 100).map(({ id, name }): RequestOptions => {
      return {
        url: `https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/RegistarServlet?subareaid=${id}`,
        userData: {
          title: name,
        },
        label: 'subarea',
        uniqueKey: createRandomUniqueKey('subarea', id),
      }
    })

    log.info(`Added ${requests.length} "subarea" requests.`)

    await crawler.addRequests(requests)

  })

}