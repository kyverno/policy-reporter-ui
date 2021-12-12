<template>
  <v-card height="350" style="height: 100%;">
    <v-card-title class="pb-4">
      Policy Types
    </v-card-title>
    <v-card-text style="min-height: 220px;">
      <wait>
        <apexchart type="donut" height="200" :options="pie.chartOptions" :series="pie.series" />
      </wait>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { ApexOptions } from 'apexcharts'
import Vue from 'vue'
import { Policy, RuleType } from '../types'
import Wait from '@/components/Wait.vue'

type Data = {
  open: boolean;
  colors: string[];
  labels: string[];
  series: number[];
}
type Computed = { pie: any }
type Props = { policies: Policy[] }
type Methods = {}

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { Wait },
  props: {
    policies: { type: Array, required: true }
  },
  data: () => ({
    open: true,
    colors: ['#089185', '#0b6ca0', '#067a11'],
    series: [],
    labels: []
  }),
  computed: {
    pie (): { series: number[]; chartOptions: ApexOptions } {
      return {
        series: this.series,
        chartOptions: {
          theme: { mode: this.$vuetify.theme.dark ? 'dark' : 'light' },
          stroke: {
            colors: this.$vuetify.theme.dark ? ['#1E1E1E'] : undefined
          },
          chart: {
            type: 'donut',
            selection: {
              enabled: false
            }
          },
          dataLabels: {
            dropShadow: {
              enabled: true,
              top: 0,
              left: 0,
              blur: 1,
              color: '#000',
              opacity: 1
            }
          },
          plotOptions: {
            pie: {
              expandOnClick: false,
              donut: {
                labels: {
                  show: true,
                  total: {
                    showAlways: true,
                    show: true
                  }
                }
              }
            }
          },
          labels: [RuleType.VALIDATION, RuleType.MUTATION, RuleType.GENERATION],
          colors: this.colors
        }
      }
    }
  },
  watch: {
    policies: {
      immediate: true,
      handler (policies: Policy[]) {
        const series = policies.reduce<number[]>((types, policy) => {
          if (policy.rules.some(r => r.type === RuleType.VALIDATION)) {
            types[0] += 1
          }
          if (policy.rules.some(r => r.type === RuleType.MUTATION)) {
            types[1] += 1
          }

          if (policy.rules.some(r => r.type === RuleType.GENERATION)) {
            types[2] += 1
          }

          return types
        }, [0, 0, 0])

        if (JSON.stringify(series) !== JSON.stringify(this.series)) {
          this.series = series
        }
      }
    }
  },
  methods: {}
})
</script>

<style scoped>
>>> code {
  padding: 16px !important;
}
</style>
