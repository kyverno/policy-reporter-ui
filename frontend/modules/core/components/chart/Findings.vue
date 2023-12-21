<template>
  <Pie :data="data" :options="(options as any)" />
</template>

<script setup lang="ts">
import { type PropType } from "vue";
import { Pie } from 'vue-chartjs';
import { type FindingCounts, Status } from '../../types'
import chroma from 'chroma-js'
import { mapStatus } from "../../mapper";
import { capilize } from "../../layouthHelper";

type Dataset = { data: number[]; label: string; backgroundColor: string[] }

const props = defineProps({
  status: { type: String as PropType<Status>, default: Status.FAIL, required: true },
  findings: { type: Object as PropType<FindingCounts>, required: true, default: () => ({ counts: [], total: 0 }) }
})

const dataset = ref<Dataset[]>([])
const labels = ref<string[]>([])
const total = ref(0)

const diff = 0.80

const colors = (status: Status, amount: number) => {
  const middle = Math.floor(amount / 2)
  const base = chroma(mapStatus(status))

  return Array.from(Array(amount).keys()).map((index) => {
    if (index < middle) {
      return base.brighten(diff * (middle - index)).hex()
    }

    return base.darken(diff * (index - middle)).hex()
  })
}

watch(() => props.findings, (findings: FindingCounts) => {
  if (!findings) return;

  const amount = findings.counts.length || 0
  const failedSet = { label: props.status, data: [], backgroundColor: colors(props.status as Status, amount) } as Dataset

  labels.value = []
  total.value = 0

  findings.counts.sort((a, b) => a.counts[props.status] - b.counts[props.status]).forEach((f) => {
    labels.value = [...labels.value, capilize(f.source)]
    failedSet.data.push(f.counts[props.status] || 0)

    total.value = total.value + (f.counts[props.status] || 0)
  })

  dataset.value = [failedSet]

}, { immediate: true })

const data = computed(() => ({
  labels: labels.value,
  datasets: dataset.value
}))

const options = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: true,
      text: `${total.value} ${capilize(props.status as Status)} Results`
    },
    legend: {
      position: 'left'
    }
  },
}))
</script>
