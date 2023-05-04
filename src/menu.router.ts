import { CheerioCrawlingContext, RouterHandler } from 'crawlee'

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

  routerRef.addHandler('menu', async ({ body, enqueueLinks }) => {

    if (!(body instanceof Buffer)) {
      throw new Error('Expected for `body` to be a Buffer.')
    }

    const menuResponse: Array<MenuResponse> = JSON.parse(body.toString())

    const subareaIdsAcc: Array<SubareaId> = []
    extractSubareaIds(menuResponse, subareaIdsAcc)

    const promises = subareaIdsAcc.slice(0, 5).map(({ id, name }) => {
      console.log('constructed', id)
      return enqueueLinks({
        urls: [`https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/RegistarServlet?subareaid=${id}`],
        userData: {
          title: name,
        },
        label: 'subarea',
      })
    })

    await Promise.all(promises)

  })

}