<template>
  <v-row>
    <v-col v-for="(filter, key) in labelFilter" :key="key" :cols="namespaced ? cols[key % 3] : 6">
      <label-filter :labels="labels[filter] || []" :name="filter" />
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import LabelFilter from './selects/LabelFilter.vue'

type Data = { labels: {[key: string]: string[]}; cols: {[key: number]: number}; interval: any }
type Props = { value: string[], source?: string, namespaced: boolean; }
type Methods = {};
type Computed = { labelFilter: string[]; refreshInterval: number };

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { LabelFilter },
  props: {
    namespaced: { type: Boolean, default: false },
    value: { type: Array, default: () => [] },
    source: { type: String, default: undefined }
  },
  data: () => ({
    interval: null,
    labels: {},
    cols: {
      0: 5,
      1: 4,
      2: 3
    }
  }),
  fetch () {
    if (this.namespaced) {
      return this.$coreAPI.namespacedReportLabels(this.source).then((labels) => {
        this.labels = labels
      })
    }
    return this.$coreAPI.clusterReportLabels(this.source).then((labels) => {
      this.labels = labels
    })
  },
  computed: mapGetters(['labelFilter', 'refreshInterval']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        if (!refreshInterval) {
          this.interval = null
          return
        }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
