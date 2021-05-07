<template>
  <v-card height="350" style="height: 100%">
      <v-card-title class="pb-4">
      Policy Categories
      </v-card-title>
    <wait>
        <v-card-text>
          <apexchart type="donut" :options="pie.chartOptions" height="200" :series="pie.series"></apexchart>
        </v-card-text>
    </wait>
  </v-card>
</template>

<script lang="ts">
import Wait from '@/components/Wait.vue';
import { ApexOptions } from 'apexcharts';
import Vue from 'vue';
import { PolicyGroups } from '../models';

type Data = { open: boolean }
type Computed = { pie: any }
type Props = { policyGroups: PolicyGroups }
type Methods = {}

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { Wait },
  props: {
    policyGroups: { type: Object, required: true },
  },
  data: () => ({
    open: true,
  }),
  computed: {
    pie(): { series: number[]; chartOptions: ApexOptions } {
      const labels = Object.keys(this.policyGroups);
      const series = Object.entries(this.policyGroups).map(([, group]) => group.length);

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
            style: {
              colors: ['#fff'],
            },
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
          labels,
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
