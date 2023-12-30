<template>
  <Bar v-bind="chart"/>
</template>

<script lang="ts" setup>
import { Bar } from 'vue-chartjs'
import { type Chart } from '../../types'
import { useStatusColors } from "~/modules/core/composables/theme";

const props = defineProps<{ data: Chart }>()

const chartColors = useChartColors()
const statusColors = useStatusColors()

const chart = computed(() => {
  return {
    style: {
      minHeight: `${ 125 + (props.data.labels.length * 25) }px`
    },
    data: {
      labels: props.data.labels,
      datasets: props.data.datasets.map(d => ({ ...d, backgroundColor: statusColors.value[d.label?.toLowerCase()] }))
    },
    options: {
      color: chartColors.value.color,
      borderColor: chartColors.value.borderColor,
      backgroundColor: chartColors.value.backgroundColor,
      height: '100%',
      indexAxis: 'y',
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        title: {
          display: true,
          text: `Resource results per Source`
        },
        legend: {
          display: true,
          position: 'bottom'
        }
      },
      scales: {
        x: {
          stacked: true,
        },
        y: {
          stacked: true
        }
      }
    }
  }
})
</script>
