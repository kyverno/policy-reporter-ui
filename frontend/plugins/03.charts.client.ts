import * as ChartJs from 'chart.js';
import { DisplayMode } from "~/modules/core/types";
import { useConfigStore } from "~/store/config";
import type { DebuggerEvent } from "vue";

const { Chart, registerables } = ChartJs;

const configColors = (mode: DisplayMode) => {
  if (mode === DisplayMode.DARK) {
    Chart.defaults.color = "#ADBABD";
    Chart.defaults.borderColor = "rgba(255,255,255,0.1)";
    Chart.defaults.backgroundColor = "rgba(255,255,0,0.1)";
    return;
  }

  Chart.defaults.color = "#666666";
  Chart.defaults.borderColor = "rgba(0,0,0,0.1)";
  Chart.defaults.backgroundColor = "rgba(0,0,0,0.1)";
}

export default defineNuxtPlugin(() => {
  const config = useConfigStore()

  Chart.register(...registerables)

  configColors(config.theme)

  config.$subscribe((mutation, state) => {
    if ((mutation.events as DebuggerEvent).key === 'displayMode') { configColors(state.displayMode) }
  })
});
