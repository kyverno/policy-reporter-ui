<template>
  <Bar v-bind="chart" />
</template>

<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import { type Chart } from '../../types'
import { useStatusColors } from "~/modules/core/composables/theme";

const props = defineProps<{ source: Chart }>()

const colors = useChartColors()
const statusColors = useStatusColors()

const chart = computed(() => {
  return {
    style: {
      minHeight: `${125 + (props.source.labels.length * 25)}px`
    },
    data: {
      labels: props.source.labels,
      datasets: props.source.datasets.map((d) => ({ ...d, backgroundColor: statusColors.value[d.label?.toLowerCase()] }))
    },
    options: {
      color: colors.value.color,
      borderColor: colors.value.borderColor,
      backgroundColor: colors.value.backgroundColor,
      height: '100%',
      indexAxis: 'y',
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        title: {
          display: true,
          text: `Results per Category`
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
