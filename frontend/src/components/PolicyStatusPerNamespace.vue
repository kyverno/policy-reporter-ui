<template>
<v-col cols="12" :md="fullWidth ? 12 : 6" v-if="!this.optional ? true : !!options.series[0].data.length">
  <v-card min-height="300" height="100%">
      <v-card-title class="pb-0">
      {{ statusText }} Policy Results per Namespace
      </v-card-title>
        <wait>
          <v-card-text class="pt-0 d-flex align-end" style="height: calc(100% - 48px)">
              <apexchart :options="options"
                        :series="options.series"
                        :height="renderHeight"
                        v-if="show"
                        style="width: 100%"
              />
          </v-card-text>
        </wait>
  </v-card>
</v-col>
</template>

<script lang="ts">
import Vue from 'vue';
import { Result, Status } from '@/models';
import { mapStatus, mapStatusText } from '@/mapper';
import Wait from './Wait.vue';

type Props = { optional: boolean; minHeight?: number; status: string; results: Result[]; fullWidth: boolean };

type Computed = {
  statusText: string;
  statusColor: string;
  renderHeight: number;
  height: number;
  options: any;
}

export default Vue.extend<{ show: boolean }, {}, Computed, Props>({
  components: { Wait },
  name: 'PolicyStatusPerNamespace',
  props: {
    optional: { default: false, type: Boolean },
    minHeight: { required: false, type: Number },
    results: { required: true, type: Array },
    status: { required: true, type: String },
    fullWidth: { default: false, type: Boolean },
  },
  data: () => ({ show: false }),
  watch: {
    height: {
      immediate: true,
      handler(height: number) {
        this.$emit('height-change', height);
      },
    },
  },
  computed: {
    renderHeight(): number {
      if (this.height > 200 && this.height > (this.minHeight || 0)) {
        return this.height;
      }

      if (this.minHeight && this.minHeight > 200) {
        return this.minHeight;
      }

      return 200;
    },
    height(): number {
      return 80 + this.options.series[0].data.length * 36;
    },
    statusText(): string {
      return mapStatusText(this.status as Status);
    },
    statusColor(): string {
      return mapStatus(this.status as Status);
    },
    options() {
      const unordnered = this.results.reduce<{ [namspace: string]: number }>((acc, item) => {
        if (!item.resource) {
          return acc;
        }

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
        colors: [this.statusColor],
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
