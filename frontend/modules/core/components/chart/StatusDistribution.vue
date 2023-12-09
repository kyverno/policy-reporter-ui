<template>
  <Pie v-bind="chart" />
</template>

<script setup lang="ts">
import type { PropType } from "vue";
import { Pie } from 'vue-chartjs';
import { type SourceFindings, Status } from '../../types'
import { mapStatus } from "../../mapper";
import { capilize } from "../../layouthHelper";

const props = defineProps({
  findings: { type: Object as PropType<SourceFindings>, required: true, default: () => ({ counts: {}, total: 0 }) }
})

const chart = computed(() => {
  if (!props.findings) return ({})

  const values = {
    [Status.PASS]: props.findings.counts[Status.PASS],
    [Status.SKIP]: props.findings.counts[Status.SKIP],
    [Status.FAIL]: props.findings.counts[Status.FAIL],
    [Status.WARN]: props.findings.counts[Status.WARN],
    [Status.ERROR]: props.findings.counts[Status.ERROR],
  }

  const data = Object.values(values).filter(c => !!c)

  const colors = Object.keys(values).filter(v => values[v]).map(s => mapStatus(s))
  const labels = Object.keys(values).filter(v => values[v]).map(s => capilize(s))

  return {
    data: {
      labels,
      datasets: [
        { data, backgroundColor: colors }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        title: {
          display: true,
          text: `${props.findings.total} ${capilize(props.findings.source || '')} Results`
        },
        legend: {
          position: 'bottom'
        }
      },
    }
  }
})
</script>
