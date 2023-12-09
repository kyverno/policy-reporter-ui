import * as ChartJs from 'chart.js';

const { Chart, registerables } = ChartJs;

export default defineNuxtPlugin(() => {
  Chart.register(...registerables)
});
