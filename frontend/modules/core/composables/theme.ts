import { useTheme } from "vuetify";
import {Severity, Status} from "../types";
import { computed } from "vue";

export const useStatusColors = () => {
  const theme = useTheme()

  const pass = computed(() => theme.current.value.colors[`status-${Status.PASS}`])
  const fail = computed(() => theme.current.value.colors[`status-${Status.FAIL}`])
  const warn = computed(() => theme.current.value.colors[`status-${Status.WARN}`])
  const error = computed(() => theme.current.value.colors[`status-${Status.ERROR}`])
  const skip = computed(() => theme.current.value.colors[`status-${Status.SKIP}`])

  return computed(() => ({
    pass: pass.value,
    fail: fail.value,
    warn: warn.value,
    error: error.value,
    skip: skip.value,
  }))
}

export const useSeverityColors = () => {
  const theme = useTheme()

  const unknown = computed(() => theme.current.value.colors[`severity-${Severity.UNKNOWN}`])
  const low = computed(() => theme.current.value.colors[`severity-${Severity.LOW}`])
  const info = computed(() => theme.current.value.colors[`severity-${Severity.INFO}`])
  const medium = computed(() => theme.current.value.colors[`severity-${Severity.MEDIUM}`])
  const high = computed(() => theme.current.value.colors[`severity-${Severity.HIGH}`])
  const critical = computed(() => theme.current.value.colors[`severity-${Severity.CRITICAL}`])

  return computed(() => ({
    unknown: unknown.value,
    low: low.value,
    info: info.value,
    medium: medium.value,
    high: high.value,
    critical: critical.value,
  }))
}
export const  useChartColors = () => {
  const theme = useTheme()
  return computed(() => {
    if (theme.current.value.dark) {
      return {
        color: '#CCCCCC',
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
