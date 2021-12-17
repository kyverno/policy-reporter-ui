import { Priority, Status } from './types'

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
  [Status.ERROR]: '#F44336'
}

const statusToDarkColor: { [status in Status]: string } = {
  [Status.SKIP]: '#424242',
  [Status.PASS]: '#1B5E20',
  [Status.WARN]: '#FF6F00',
  [Status.FAIL]: '#D32F2F',
  [Status.ERROR]: '#B71C1C'
}

const statusToText: { [status in Status]: string } = {
  [Status.SKIP]: 'Skipped',
  [Status.PASS]: 'Passing',
  [Status.WARN]: 'Warning',
  [Status.FAIL]: 'Failing',
  [Status.ERROR]: 'Errored'
}

export const mapPriority = (priority: Priority): string => priorityToColor[priority] || priorityToColor[Priority.DEBUG]
export const mapStatus = (status: Status): string => statusToColor[status] || statusToColor[Status.SKIP]
export const mapStatusText = (status: Status): string => statusToText[status] || statusToText[Status.SKIP]
export const mapDarkStatus = (status: Status): string => statusToDarkColor[status] || statusToDarkColor[Status.SKIP]
