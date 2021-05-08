<template>
  <v-card height="350" style="height: 100%">
      <v-card-title class="pb-4">
      Policy Categories
      </v-card-title>
    <wait>
        <v-card-text>
          <apexchart type="donut" height="200" :options="pie.chartOptions" :series="pie.series"></apexchart>
        </v-card-text>
    </wait>
  </v-card>
</template>

<script lang="ts">
import Wait from '@/components/Wait.vue';
import { ApexOptions } from 'apexcharts';
import Vue from 'vue';
import chroma from 'chroma-js';
import { PolicyGroups } from '../models';

type Data = {
  open: boolean;
  colors: string[];
  labels: string[];
  series: number[];
}

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
    colors: chroma.scale(['#065684', '#089185', '#067a11', '#097a15']).mode('lch').colors(9),
    series: [],
    labels: [],
  }),
  watch: {
    policyGroups: {
      immediate: true,
      handler(policyGroups: PolicyGroups) {
        const labels = Object.keys(policyGroups);
        if (JSON.stringify(labels) !== JSON.stringify(this.labels)) {
          this.labels = labels;
        }

        const series = Object.entries(policyGroups).map(([, group]) => group.length);
        if (JSON.stringify(series) !== JSON.stringify(this.series)) {
          this.series = series;
        }
      },
    },
  },
  computed: {
    pie(): { series: number[]; chartOptions: ApexOptions } {
      return {
        series: this.series,
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
          labels: this.labels,
          colors: this.colors,
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
