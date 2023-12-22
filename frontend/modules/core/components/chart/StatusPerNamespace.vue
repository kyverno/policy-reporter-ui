<template>
  <Bar v-bind="chart" />
</template>

<script setup lang="ts">
import { Bar } from 'vue-chartjs'
import { type Filter, type NamespacedStatusCount, Status } from '../../types'
import { capilize } from "../../layouthHelper"
import { mapStatus } from '../../mapper'
import { NamespacedKinds, ResourceFilter } from "~/modules/core/provider/dashboard";
import type { Ref } from "vue";

const props = defineProps<{ source: string }>()

const filter = inject<Ref<Filter>>(ResourceFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(NamespacedKinds, ref<string[]>([]))

const { data, refresh } = useAPI<NamespacedStatusCount[]>(
    (api) => api.namespacedStatusCount({
      ...filter.value,
      sources: [props.source],
      kinds: kinds.value.length ? kinds.value : undefined
    }), {
      default: () => [],
    }
);

watch(kinds, () => refresh())

const chart = computed(() => {
  if (!data.value) return ({})

  const list: { [key: string]: { [key in Status]: number }} = {}

  data.value.forEach(f => {
    f.items.forEach(i => {
      if (!list[i.namespace]) {
        list[i.namespace] = {
          [Status.PASS]: 0,
          [Status.SKIP]: 0,
          [Status.FAIL]: 0,
          [Status.WARN]: 0,
          [Status.ERROR]: 0,
        }
      }

      list[i.namespace][f.status] = i.count
    })
  })

  const ordered: any = Object.keys(list).sort((a,b) => a.localeCompare(b)).reduce((obj, k) => ({
    ...obj,
    [k]: list[k]
  }), {})

  const labels = Object.keys(ordered)

  const sets: { [key in Omit<Status, Status.SKIP>]: { data: number[]; label: string; backgroundColor: string } } = {
    [Status.PASS]: { data: [], label: capilize(Status.PASS), backgroundColor: mapStatus(Status.PASS)},
    [Status.FAIL]: { data: [], label: capilize(Status.FAIL), backgroundColor: mapStatus(Status.FAIL)},
    [Status.WARN]: { data: [], label: capilize(Status.WARN), backgroundColor: mapStatus(Status.WARN)},
    [Status.ERROR]: { data: [], label: capilize(Status.ERROR), backgroundColor: mapStatus(Status.ERROR)},
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
      height: '100%',
      indexAxis: 'y',
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        title: {
          display: true,
          text: `${capilize(props.source)} Results per Namespace`
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
