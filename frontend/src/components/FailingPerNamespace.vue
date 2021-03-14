<template>
<v-card min-height="300">
    <v-card-title class="pb-0">
    Failing Policies per Namespace
    </v-card-title>
    <v-card-text class="pt-0">
    <apexchart :options="options" :series="options.series" :height="height" v-if="show" />
    </v-card-text>
</v-card>
</template>

<script lang="ts">
import Vue from 'vue';
import { NamespacePolicyMap } from '@/models';

export default Vue.extend<{ show: boolean }, {}, { height: number; options: any }, { reports: NamespacePolicyMap }>({
  name: 'FailingPerNamespace',
  props: {
    reports: { required: true, type: Object },
  },
  data: () => ({ show: false }),
  computed: {
    height(): number {
      const height = 80 + this.options.series[0].data.length * 36;

      if (height < 300) {
        return 300;
      }

      return height;
    },
    options() {
      const result = Object.keys(this.reports).reduce<{ [namspace: string]: number }>((acc, item) => {
        if (this.reports[item].summary.fail === 0) {
          return acc;
        }

        acc[item] = this.reports[item].summary.fail;

        return acc;
      }, {});

      const data = Object.values(result);
      const categories = Object.keys(result);

      return {
        series: [{
          data,
        }],
        chart: {
          type: 'bar',
          toolbar: {
            show: false,
          },
        },
        plotOptions: {
          bar: {
            horizontal: true,
            dataLabels: {
              position: 'bottom',
            },
          },
        },
        dataLabels: {
          enabled: true,
          offsetX: 30,
        },
        xaxis: {
          categories,
          min: 0,
        },
        tooltip: {
          theme: 'dark',
          x: {
            show: false,
          },
          y: {
            title: {
              formatter(_: string, config: { dataPointIndex: number }) {
                return categories[config.dataPointIndex] || '';
              },
            },
          },
        },
      };
    },
  },
  created() {
    setTimeout(() => { this.show = true; }, 500);
  },
});
</script>
