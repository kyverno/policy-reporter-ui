<template>
  <Bar v-bind="chart" />
</template>

<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import { type ResourceStatusCount, Status } from '../../types'
import { capilize } from "../../layouthHelper"
import { mapStatus } from '../../mapper'
import type { PropType } from "vue";

const props = defineProps<{ data: ResourceStatusCount[] }>()

const chart = computed(() => {

  const list: { [key: string]: { [key in Status]: number }} = {}

  props.data?.forEach(f => {
    f.items.forEach(i => {
      if (!list[i.source]) {
        list[i.source] = {
          [Status.PASS]: 0,
          [Status.SKIP]: 0,
          [Status.FAIL]: 0,
          [Status.WARN]: 0,
          [Status.ERROR]: 0,
        }
      }

      list[i.source][f.status] = i.count
    })
  })

  const ordered: any = Object.keys(list).sort((a,b) => a.localeCompare(b)).reduce((obj, k) => ({
    ...obj,
    [k]: list[k]
  }), {})

  const sources = Object.keys(ordered)

  const sets: { [key in Omit<Status, Status.SKIP>]: { data: number[]; label: string; backgroundColor: string } } = {
    [Status.PASS]: { data: [], label: capilize(Status.PASS), backgroundColor: mapStatus(Status.PASS)},
    [Status.FAIL]: { data: [], label: capilize(Status.FAIL), backgroundColor: mapStatus(Status.FAIL)},
    [Status.WARN]: { data: [], label: capilize(Status.WARN), backgroundColor: mapStatus(Status.WARN)},
    [Status.ERROR]: { data: [], label: capilize(Status.ERROR), backgroundColor: mapStatus(Status.ERROR)},
  }

  sources.forEach((source) => {
    sets[Status.PASS].data.push(ordered[source][Status.PASS])
    sets[Status.FAIL].data.push(ordered[source][Status.FAIL])
    sets[Status.WARN].data.push(ordered[source][Status.WARN])
    sets[Status.ERROR].data.push(ordered[source][Status.ERROR])
  })

  return {
    style: {
      minHeight: `${125 + (sources.length * 25)}px`
    },
    data: {
      labels: sources.map(l => capilize(l)),
      datasets: Object.values(sets)
    },
    options: {
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
