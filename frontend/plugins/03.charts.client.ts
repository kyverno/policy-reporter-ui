import * as ChartJs from 'chart.js';
import { DisplayMode } from "~/modules/core/types";
import { useConfigStore } from "~/store/config";

const { Chart, registerables } = ChartJs;

const configColors = (mode: DisplayMode) => {
  if ([DisplayMode.DARK, DisplayMode.COLOR_BLIND_DARK].includes(mode)) {
    Chart.defaults.color = "#CCCCCC";
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
    configColors(state.displayMode)
  })
});
