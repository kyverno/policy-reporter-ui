import { useTheme } from "vuetify";
import { mapDarkStatus, mapStatus } from "../mapper";
import { Severity, Status } from "../types";
import { computed } from "vue";

const statusToBtnColor: { [status in Status]: string } = {
  [Status.SKIP]: '#EEEEEE',
  [Status.PASS]: '#00E676',
  [Status.WARN]: '#FF6D00',
  [Status.FAIL]: '#FF5252',
  [Status.ERROR]: '#E53935'
}

export const useStatusMapper = (status: Status) => {
  const theme = useTheme()

  return computed(() => theme.current.value.dark ? mapDarkStatus(status) : mapStatus(status))
}

export const useStatusBtnColor = (status: Status) => {
  const theme = useTheme()

  return computed(() => theme.current.value.dark ? statusToBtnColor[status] : mapStatus(status))
}

export const useStatusColors = () => {
  const pass = useStatusMapper(Status.PASS)
  const fail = useStatusMapper(Status.FAIL)
  const warn = useStatusMapper(Status.WARN)
  const error = useStatusMapper(Status.ERROR)
  const skip = useStatusMapper(Status.SKIP)

  return computed(() => ({
    pass: pass.value,
    fail: fail.value,
    warn: warn.value,
    error: error.value,
    skip: skip.value,
  }))
}
export const  useChartColors = () => {
  const theme = useTheme()
  return computed(() => {
    if (theme.current.value.dark) {
      return {
        color: '#ADBABD',
        borderColor: "rgba(255,255,255,0.1)",
        backgroundColor: "rgba(255,255,255,0.1)",
        element: "rgba(255,255,255,0.4)",
      }
    }

    return {
      color: '#666666',
      borderColor: "rgba(0,0,0,0.1)",
      backgroundColor: "rgba(0,0,0,0.1)",
      element: "rgba(0,0,0,0.4)",
    }
  })
}

export const useBGColor = () => {
  const theme = useTheme()

  return computed(() => {
    if (theme.current.value.dark) {
      return 'bg-black'
    }

    return 'bg-grey-lighten-3'
  })
}

const severityToColor = {
  [Severity.INFO]: 'grey',
  [Severity.LOW]: 'info',
  [Severity.MEDIUM]: 'warning',
  [Severity.HIGH]: 'red-darken-2',
  [Severity.CRITICAL]: 'red-accent-3',
}

const severityToDarkColor = {
  [Severity.INFO]: 'light-blue-accent-1',
  [Severity.LOW]: 'blue-accent-1',
  [Severity.MEDIUM]: 'amber-accent-3',
  [Severity.HIGH]: 'red-accent-1',
  [Severity.CRITICAL]: 'red-accent-2',
}

export const useSeverityColor = (severity: Severity) => {
  const theme = useTheme()

  return computed(() => theme.current.value.dark ? severityToDarkColor[severity] : severityToColor[severity])
}
