<template>
  <Bar v-bind="chart" />
</template>

<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import { type Source, Status } from '../../types'
import { capilize } from "../../layouthHelper"
import { useStatusColors } from "~/modules/core/composables/theme";

const props = defineProps<{ source: Source }>()

const colors = useChartColors()
const statusColors = useStatusColors()

const chart = computed(() => {
  const list: { [key: string]: { [key in Status]: number }} = {}

  props.source.categories.forEach(f => {
    list[f.name || 'other'] = {
      [Status.PASS]: f.pass,
      [Status.SKIP]: f.skip,
      [Status.FAIL]: f.fail,
      [Status.WARN]: f.warn,
      [Status.ERROR]: f.error,
    }
  })

  const ordered: any = Object.keys(list).sort((a,b) => a.localeCompare(b)).reduce((obj, k) => ({
    ...obj,
    [k]: list[k]
  }), {})

  const labels = Object.keys(ordered)

  const sets: { [key in Omit<Status, Status.SKIP>]: { data: number[]; label: string; backgroundColor: string } } = {
    [Status.PASS]: { data: [], label: capilize(Status.PASS), backgroundColor: statusColors.value.pass },
    [Status.FAIL]: { data: [], label: capilize(Status.FAIL), backgroundColor: statusColors.value.fail },
    [Status.WARN]: { data: [], label: capilize(Status.WARN), backgroundColor: statusColors.value.warn },
    [Status.ERROR]: { data: [], label: capilize(Status.ERROR), backgroundColor: statusColors.value.error },
  }

  labels.forEach((ns) => {
    sets[Status.PASS].data.push(ordered[ns][Status.PASS])
    sets[Status.FAIL].data.push(ordered[ns][Status.FAIL])
    sets[Status.WARN].data.push(ordered[ns][Status.WARN])
    sets[Status.ERROR].data.push(ordered[ns][Status.ERROR])
  })

  return {
    style: {
      minHeight: `${125 + (labels.length * 25)}px`
    },
    data: {
      labels,
      datasets: Object.values(sets)
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
