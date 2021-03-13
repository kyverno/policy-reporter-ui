<template>
<v-col cols="6" v-if="!this.optional ? true : !!options.series[0].data.length">
  <v-card min-height="300" height="100%">
      <v-card-title class="pb-0">
      {{ statusText }} Policy Results per Namespace
      </v-card-title>
      <v-card-text class="pt-0 d-flex align-end" style="height: calc(100% - 48px)">
      <apexchart :options="options" :series="options.series" :height="height" v-if="show" style="width: 100%" />
      </v-card-text>
  </v-card>
</v-col>
</template>

<script lang="ts">
import Vue from 'vue';
import { Result } from '@/models';

type Props = { optional: boolean; minHeight?: number; statusText: string; results: Result[] };

export default Vue.extend<{ show: boolean }, {}, { height: number; options: any }, Props>({
  name: 'PolicyStatusPerNamespace',
  props: {
    optional: { default: false, type: Boolean },
    minHeight: { required: false, type: Number },
    results: { required: true, type: Array },
    statusText: { required: true, type: String },
  },
  data: () => ({ show: false }),
  watch: {
    height(height: number) {
      if (!height) return;

      this.$emit('height-change', height);
    },
  },
  computed: {
    height(): number {
      const height = 80 + this.options.series[0].data.length * 36;

      if (this.minHeight && height < this.minHeight) {
        return this.minHeight;
      }

      if (height < 200) {
        return 200;
      }

      return height;
    },
    options() {
      const unordnered = this.results.reduce<{ [namspace: string]: number }>((acc, item) => {
        if (!item.resource.namespace) {
          return acc;
        }

        acc[item.resource.namespace] = (acc[item.resource.namespace] || 0) + 1;

        return acc;
      }, {});

      const ordered = Object.keys(unordnered).sort().reduce<{ [namspace: string]: number }>((acc, key) => {
        acc[key] = unordnered[key];

        return acc;
      }, {});

      const data = Object.values(ordered);
      const categories = Object.keys(ordered);

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
