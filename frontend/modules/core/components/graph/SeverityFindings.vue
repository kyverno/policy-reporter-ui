<template>
  <wait :time="time">
    <Pie :data="data" :options="(options as any)" />
    <template #placeholder>
      <v-container fluid>
        <v-row>
          <v-col class="justify-center align-center text-center d-flex">
            <v-progress-circular indeterminate size="268" width="10" color="primary" />
          </v-col>
        </v-row>
      </v-container>
    </template>
  </wait>
</template>

<script setup lang="ts">
import { type PropType } from "vue";
import { Pie } from 'vue-chartjs';
import { type Findings, Severity } from '../../types'
import chroma from 'chroma-js'
import { capilize } from "../../layouthHelper";
import {useTheme} from "vuetify";

const props = defineProps({
  severity: { type: String as PropType<Severity>, default: Severity.MEDIUM, required: true },
  time: { type: Number, default: 400, required: false },
  data: { type: Object as PropType<Findings>, required: true, default: (): Findings => ({ labels: [], datasets: [], name: "0" }) }
})

const diff = 0.80
const theme = useTheme()

const colors = (severity: Severity, amount: number) => {
  const middle = Math.floor(amount / 2)
  const base = chroma(theme.current.value.colors[`severity-${severity}`])

  return Array.from(Array(amount).keys()).map((index) => {
    if (index < middle) {
      return base.brighten(diff * (middle - index)).hex()
    }

    return base.darken(diff * (index - middle)).hex()
  })
}

const chartColors = useChartColors()
const data = computed(() => ({
  labels: props.data?.labels,
  datasets: props.data?.datasets.map(d => ({ ...d, backgroundColor: colors(props.severity as Severity, d.data.length) }))
}))

const options = computed(() => ({
  color: chartColors.value.color,
  borderColor: chartColors.value.borderColor,
  backgroundColor: chartColors.value.backgroundColor,
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: true,
      text: `${props.data?.name} ${capilize(props.severity as Severity)} Results`
    },
    legend: {
      position: 'left'
    }
  },
}))
</script>
