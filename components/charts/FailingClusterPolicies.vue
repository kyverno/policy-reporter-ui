<template>
  <v-card full-height style="height: 100%;" min-height="300">
    <v-card-title class="pb-0">
      Failing Cluster Policies
    </v-card-title>
    <v-card-text :style="`height: calc(100% - 48px); font-size: ${size}rem !important;`" class="text-center text-h1 d-flex justify-center align-center">
      <span v-if="!waiting">{{ counter }}</span>
      <v-progress-circular v-else indeterminate size="100" />
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { Status } from '~/policy-reporter-plugins/core/types'

export default Vue.extend<{ counter: number; interval: any; waiting: boolean }, {}, { size: number }, {}>({
  name: 'FailingClusterPolicies',
  data: () => ({
    counter: 0,
    interval: null,
    waiting: true
  }),
  async fetch () {
    const statusCount = await this.$coreAPI.statusCount({ status: [Status.FAIL] })

    const failed = statusCount[0] || { status: Status.FAIL, count: 0 }

    this.counter = failed.count
  },
  computed: {
    size () {
      const counterLength = this.counter.toString().length

      if (counterLength <= 3) { return 12 }

      return 12 - ((counterLength - 2))
    },
    ...mapGetters(['refreshInterval'])
  },
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  created () {
    setTimeout(() => { this.waiting = false }, 500)
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
