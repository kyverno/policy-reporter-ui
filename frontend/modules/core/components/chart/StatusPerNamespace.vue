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
import { execOnChange } from "~/helper/compare";
import { useStatusColors } from "~/modules/core/composables/theme";

const props = defineProps<{ title?: string; source: string }>()

const filter = inject<Ref<Filter>>(ResourceFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(NamespacedKinds, ref<string[]>([]))

const computedFilter = computed(() => ({
  ...filter.value,
  sources: [props.source],
  kinds: kinds.value.length ? kinds.value : undefined
}))

const { data, refresh } = useAPI(
    (api) => api.namespacedStatusCount(computedFilter.value), {
      default: () => ({}),
    }
);

watch(computedFilter, (n,o) => execOnChange(n,o, () => refresh()))

const colors = useChartColors()
const statusColors = useStatusColors()

const chart = computed(() => {
  if (!data.value) return ({})

  const list: { [key: string]: { [key in Status]: number }} = {}

  const ordered: any = Object.keys(data.value).sort((a,b) => a.localeCompare(b)).reduce((obj, k) => ({
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
          text: `${props.title ?? capilize(props.source)} Results per Namespace`
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
