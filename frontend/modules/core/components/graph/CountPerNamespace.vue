<template>
  <Bar v-bind="chart" />
</template>

<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import { capilize } from "../../layouthHelper"
import {useChartColors, useSeverityColors, useStatusColors} from "~/modules/core/composables/theme";
import type { Chart } from "~/modules/core/types";

const props = defineProps<{ title?: string; data: Chart }>()
const colors = useChartColors()
const statusColors = useStatusColors()
const severityColors = useSeverityColors()

const config = computed(() => {
  if (props.data.type === 'severity') {
    return {
      colors: severityColors.value,
      title: `Severities per Namespace`
    }
  }

  return {
    colors: statusColors.value,
    title: `Results per Namespace`
  }
})

const chart = computed(() => {
  return {
    style: {
      minHeight: `${125 + (props.data.labels.length * 25)}px`
    },
    data: {
      labels: props.data.labels,
      datasets: props.data.datasets.map((d) => ({ ...d, backgroundColor: config.value.colors[d.label?.toLowerCase()] }))
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
          text: `${props.title ?? capilize(props.data.name )} ${config.value.title}`
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
