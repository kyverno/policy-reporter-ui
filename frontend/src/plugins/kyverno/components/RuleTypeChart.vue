<template>
  <v-card height="350" style="height: 100%">
      <v-card-title class="pb-4">
      Policy Types
      </v-card-title>
    <wait>
        <v-card-text>
        <apexchart type="donut" height="200"  :options="pie.chartOptions" :series="pie.series"></apexchart>
        </v-card-text>
    </wait>
  </v-card>
</template>

<script lang="ts">
import Wait from '@/components/Wait.vue';
import { ApexOptions } from 'apexcharts';
import Vue from 'vue';
import { Policy, RuleType } from '../models';

type Data = { open: boolean; search: string; expanded: string[] }
type Computed = { pie: any }
type Props = { policies: Policy[] }
type Methods = {}

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { Wait },
  props: {
    policies: { type: Array, required: true },
  },
  data: () => ({ open: true, search: '', expanded: [] }),
  computed: {
    pie(): { series: number[]; chartOptions: ApexOptions } {
      const series = this.policies.reduce<number[]>((types, policy) => {
        if (policy.rules.some((r) => r.type === RuleType.VALIDATION)) {
          types[0] += 1;
        }
        if (policy.rules.some((r) => r.type === RuleType.MUTATION)) {
          types[1] += 1;
        }

        if (policy.rules.some((r) => r.type === RuleType.GENERATION)) {
          types[2] += 1;
        }

        return types;
      }, [0, 0, 0]);

      return {
        series,
        chartOptions: {
          chart: {
            type: 'donut',
            selection: {
              enabled: false,
            },
          },
          dataLabels: {
            dropShadow: {
              enabled: true,
              top: 0,
              left: 0,
              blur: 1,
              color: '#000',
              opacity: 1,
            },
          },
          plotOptions: {
            pie: {
              expandOnClick: false,
              donut: {
                labels: {
                  show: true,
                  total: {
                    showAlways: true,
                    show: true,
                  },
                },
              },
            },
          },
          labels: [RuleType.VALIDATION, RuleType.MUTATION, RuleType.GENERATION],
        },
      };
    },
  },
  methods: {},
});
</script>

<style scoped>
>>> code {
  padding: 16px!important;
}
</style>
