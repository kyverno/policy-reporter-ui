<template>
  <loader :loading="loading" :error="$fetchState.error">
    <v-container fluid class="pt-6 px-6">
      <v-row>
        <policy-status-per-namespace
          v-if="config.policyReports"
          :values="failedCounters"
          status="fail"
          xs-cols="12"
          :md-cols="config.clusterPolicyReports ? 8 : 12"
        />
        <v-col
          v-if="config.clusterPolicyReports"
          xs-cols="12"
          :md-cols="config.policyReports ? 4 : 12"
        >
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
      <cluster-policy-report-table v-if="config.clusterPolicyReports" status="fail" :filter="{ sources: [source] }" />
      <policy-report-table v-if="config.policyReports" status="fail" :filter="{ sources: [source] }" />
    </v-container>
  </loader>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import PolicyStatusPerNamespace from '~/components/charts/PolicyStatusPerNamespace.vue'
import { Cluster, DashboardConfig, Status } from '~/policy-reporter-plugins/core/types'

type FailCounters = { namespaces: string[]; counts: number[] }

type Data = {
  interval: any;
  sources: string[];
  loading: boolean;
  failedCounters: FailCounters
}

type Computed = {
  refreshInterval: number;
  config: DashboardConfig;
  currentCluster?: Cluster;
}

export default Vue.extend<Data, {}, Computed, {}>({
  components: { PolicyStatusPerNamespace },
  data: () => ({
    interval: null,
    sources: [],
    loading: true,
    failedCounters: { namespaces: [], counts: [] }
  }),
  async fetch () {
    const sourceAPIs: Promise<string[]>[] = []

    if (this.config.policyReports) {
      sourceAPIs.push(this.$coreAPI.namespacedSources())
    }
    if (this.config.clusterPolicyReports) {
      sourceAPIs.push(this.$coreAPI.clusterSources())
    }

    const mixedSources = await Promise.all(sourceAPIs)

    this.sources = [...new Set<string>(mixedSources.reduce<string[]>((acc, source) => [...acc, ...source], []))]

    if (this.config.policyReports) {
      const statusCount = await this.$coreAPI.namespacedStatusCount({ status: [Status.FAIL] })

      this.failedCounters = statusCount[0].items.reduce<{ namespaces: string[]; counts: number[] }>((acc, statusCount) => {
        acc.namespaces.push(statusCount.namespace)
        acc.counts.push(statusCount.count)

        return acc
      }, { namespaces: [], counts: [] })
    }

    this.loading = false
  },
  computed: mapGetters({ refreshInterval: 'refreshInterval', config: 'dashboardConfig', cluster: 'currentCluster' }),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    },
    cluster () {
      this.loading = true
      this.$fetch()
    }
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
