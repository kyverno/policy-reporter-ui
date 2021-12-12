<template>
  <loader :loading="loading" :error="$fetchState.error">
    <v-container fluid class="pt-6 px-6">
      <v-row>
        <policy-status-per-namespace :values="failedCounters" status="fail" :md-cols="8" />
        <v-col cols="12" md="4">
          <failing-cluster-policies />
        </v-col>
      </v-row>
    </v-container>
    <v-container v-for="source in sources" :key="source" fluid class="px-6">
      <v-row>
        <v-col>
          <v-card>
            <v-card-title class="text-h3">
              {{ source }}
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>
      <cluster-policy-report-table status="fail" :filter="{ sources: [source] }" />
      <policy-report-table status="fail" :filter="{ sources: [source] }" />
    </v-container>
  </loader>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import PolicyStatusPerNamespace from '~/components/charts/PolicyStatusPerNamespace.vue'
import { Status } from '~/policy-reporter-plugins/core/types'

export default Vue.extend({
  components: { PolicyStatusPerNamespace },
  data: () => ({
    interval: null as any,
    sources: [] as string[],
    loading: true as boolean,
    failedCounters: { namespaces: [], counts: [] } as { namespaces: string[]; counts: number[] }
  }),
  async fetch () {
    const [clusterSources, namespacedSources] = await Promise.all([
      this.$coreAPI.clusterSources(),
      this.$coreAPI.namespacedSources()
    ])

    this.sources = [...new Set<string>([...clusterSources, ...namespacedSources])]

    const statusCount = await this.$coreAPI.namespacedStatusCount({ status: [Status.FAIL] })

    this.failedCounters = statusCount[0].items.reduce<{ namespaces: string[]; counts: number[] }>((acc, statusCount) => {
      acc.namespaces.push(statusCount.namespace)
      acc.counts.push(statusCount.count)

      return acc
    }, { namespaces: [], counts: [] })

    this.loading = false
  },
  computed: mapGetters(['refreshInterval']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
