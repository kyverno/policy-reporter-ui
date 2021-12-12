<template>
  <v-col v-if="!optional ? true : !!values.namespaces.length" cols="12" :md="fullWidth ? 12 : mdCols">
    <v-card min-height="300" height="100%">
      <v-card-title class="pb-0">
        {{ statusText }} Policy Results per Namespace
      </v-card-title>
      <v-card-text class="pt-0 d-flex align-end" style="height: calc(100% - 48px);">
        <apexchart
          :options="options"
          :series="options.series"
          :height="renderHeight"
          style="width: 100%;"
        />
      </v-card-text>
    </v-card>
  </v-col>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapStatus, mapStatusText } from '~/policy-reporter-plugins/core/mapper'
import { Status } from '~/policy-reporter-plugins/core/types'

type Props = {
  optional: boolean;
  minHeight?: number;
  status: Status;
  values: { counts: number[]; namespaces: string[]; };
  fullWidth: boolean
  mdCols: number
};

type Computed = {
  statusText: string;
  statusColor: string;
  renderHeight: number;
  height: number;
  options: any;
}

type Data = {}

export default Vue.extend<Data, {}, Computed, Props>({
  name: 'PolicyStatusPerNamespace',
  props: {
    optional: { default: false, type: Boolean },
    minHeight: { required: false, type: Number, default: 0 },
    values: { required: true, type: Object },
    status: { required: true, type: String as Vue.PropType<Status> },
    fullWidth: { default: false, type: Boolean },
    mdCols: { default: 6, type: Number }
  },
  computed: {
    renderHeight (): number {
      if (this.height > 200 && this.height > (this.minHeight || 0)) {
        return this.height
      }

      if (this.minHeight && this.minHeight > 200) {
        return this.minHeight
      }

      return 200
    },
    height (): number {
      return 80 + this.options.series[0].data.length * 36
    },
    statusText (): string {
      return mapStatusText(this.status)
    },
    statusColor (): string {
      return mapStatus(this.status)
    },
    options () {
      return {
        theme: { mode: this.$vuetify.theme.dark ? 'dark' : 'light' },
        series: [{
          data: this.values.counts
        }],
        chart: {
          type: 'bar',
          toolbar: {
            show: false
          }
        },
        colors: [this.statusColor],
        plotOptions: {
          bar: {
            horizontal: true,
            dataLabels: {
              position: 'bottom'
            }
          }
        },
        dataLabels: {
          enabled: true,
          offsetX: 30
        },
        xaxis: {
          categories: this.values.namespaces,
          min: 0
        },
        tooltip: {
          theme: 'dark',
          x: {
            show: false
          },
          y: {
            title: {
              formatter: (_: string, config: { dataPointIndex: number }) => {
                return this.values.namespaces[config.dataPointIndex] || ''
              }
            }
          }
        }
      }
    }
  },
  watch: {
    height: {
      immediate: true,
      handler (height: number) {
        this.$emit('height-change', height)
      }
    }
  }
})
</script>
