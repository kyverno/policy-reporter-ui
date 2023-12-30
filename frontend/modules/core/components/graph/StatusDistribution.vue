<template>
  <Pie v-bind="chart" />
</template>

<script setup lang="ts">
import type { PropType } from "vue";
import { Pie } from 'vue-chartjs';
import { capilize } from "../../layouthHelper";
import { useStatusColors } from "~/modules/core/composables/theme";
import type { Chart } from "~/modules/core/types";

const props = defineProps({
  data: { type: Object as PropType<Chart>, required: true, default: () => ({ labels: [], datasets: [], name: "" }) },
  title: { type: String, required: false },
})

const chartColors = useChartColors()
const statusColors = useStatusColors()

const chart = computed(() => {
  const colors = props.data!.labels.map(s => statusColors.value[s.toLowerCase()])
  const total: number = props.data!.datasets[0].data.reduce((sum: number, i: number) => sum + i, 0)

  return {
    data: {
      labels: props.data?.labels,
      datasets: [
        { data: props.data?.datasets[0].data, backgroundColor: colors }
      ]
    },
    options: {
      color: chartColors.value.color,
      borderColor: chartColors.value.borderColor,
      backgroundColor: chartColors.value.backgroundColor,
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        title: {
          display: true,
          text: `${total} ${props.title ?? capilize(props.data!.name || '')} Results`
        },
        legend: {
          position: 'bottom'
        }
      },
    }
  }
})
</script>
