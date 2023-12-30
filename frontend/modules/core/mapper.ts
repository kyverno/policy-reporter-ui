import { sortByKeys } from './layouthHelper'
import { Priority, Status, type MappedResult, type Dictionary, type ListResult } from './types'

const priorityToColor: { [key in Priority]: string } = {
  [Priority.SUCCESS]: 'green lighten-2',
  [Priority.DEBUG]: 'light-blue lighten-2',
  [Priority.INFO]: 'green darken-1',
  [Priority.WARNING]: 'orange lighten-1',
  [Priority.ERROR]: 'red darken-3',
  [Priority.CRITICAL]: 'red darken-4'
}

const statusToColor: { [status in Status]: string } = {
  [Status.SKIP]: '#E0E0E0',
  [Status.PASS]: '#43A047',
  [Status.WARN]: '#FB8C00',
  [Status.FAIL]: '#EF5350',
  [Status.ERROR]: '#950011'
}

const statusToDarkColor: { [status in Status]: string } = {
  [Status.SKIP]: '#424242',
  [Status.PASS]: '#1B5E20',
  [Status.WARN]: '#FF6F00',
  [Status.FAIL]: '#D32F2F',
  [Status.ERROR]: '#950011'
}

export const mapPriority = (priority: Priority): string => priorityToColor[priority] || priorityToColor[Priority.DEBUG]
export const mapStatus = (status: Status): string => statusToColor[status] || statusToColor[Status.SKIP]
export const mapDarkStatus = (status: Status): string => statusToDarkColor[status] || statusToDarkColor[Status.SKIP]

const maxChipLength = 75

export const mapResults = ({ items, count }: { items: ListResult[], count: number }): { results: MappedResult[], count: number } => {
  const results: MappedResult[] = items.map(({ properties, ...result }) => {
    const chips: Dictionary = {}
    const cards: Dictionary = {}
    let hasProps: boolean = false

    for (const prop in properties) {
      if (prop == 'resultID') { continue }

      if (properties[prop].length > maxChipLength) {
        cards[prop] = properties[prop]
      } else {
        chips[prop] = properties[prop]
      }
      hasProps = true
    }

    return {
      ...result,
      properties: {},
      cards: sortByKeys(cards),
      chips: sortByKeys(chips),
      hasProps
    }
  })

  return { results, count }
}
