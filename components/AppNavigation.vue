<template>
  <v-navigation-drawer :value="value" app clipped width="305" @input="$emit('input', $event)">
    <v-list>
      <app-navigation-item icon="mdi-view-dashboard" route="/" title="Dashboard" />

      <template v-if="views.policyReports">
        <app-navigation-item
          v-if="namespacedPages.length <= 1"
          icon="mdi-file-chart"
          :route="(namespacedPages[0] || { path: '/policy-reports' }).path"
          title="Policy Reports"
        />
        <app-navigation-group
          v-else
          icon="mdi-file-chart"
          title="Policy Reports"
          :items="namespacedPages"
        />
      </template>

      <template v-if="views.clusterPolicyReports">
        <app-navigation-item
          v-if="clusterPages.length <= 1"
          icon="mdi-file-chart"
          :route="(clusterPages[0] || { path: '/cluster-policy-reports' }).path"
          title="ClusterPolicy Reports"
        />
        <app-navigation-group
          v-else
          icon="mdi-file-chart"
          title="ClusterPolicy Reports"
          :items="clusterPages"
        />
      </template>

      <app-navigation-item v-if="showLogs" icon="mdi-console" route="/logs" title="Logs" />

      <template v-if="views.kyvernoPolicies">
        <app-navigation-item v-if="showPlugin('kyverno')" route="/kyverno-plugin" title="Kyverno Policies" exact>
          <template #icon>
            <v-list-item-icon>
              <kyverno-icon style="height: 24px; width: 24px;" />
            </v-list-item-icon>
          </template>
        </app-navigation-item>
      </template>

      <template v-if="views.kyvernoVerifyImages">
        <app-navigation-item
          v-if="showPlugin('kyverno')"
          route="/kyverno-plugin/verify-image-rules"
          title="Kyverno VerifyImages"
          icon="mdi-shield-check"
        />
      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import AppNavigationItem from './AppNavigationItem.vue'
import { Cluster, ViewsCofig } from '~/policy-reporter-plugins/core/types'

type Page = { title: string; path: string }

type Data = {
  interval: any,
  clusterPages: Page[];
  namespacedPages: Page[];
}

type Props = {
  value: boolean;
  plugins: string[];
  views: ViewsCofig;
}

type Computed = {
  refreshInterval: number;
  currentCluster?: Cluster
  showLogs: boolean;
}

type Methdos = {
  showPlugin(plugin: string): boolean;
}

type Plugin = 'kyverno'

export default Vue.extend<Data, Methdos, Computed, Props>({
  components: { AppNavigationItem },
  props: {
    value: { type: Boolean, default: false },
    plugins: { type: Array, default: () => [] },
    views: { type: Object, required: true }
  },
  data: () => ({
    interval: null,
    clusterPages: [],
    namespacedPages: []
  }),
  async fetch () {
    const [clusterSources, namespacedSources] = await Promise.all([
      this.$coreAPI.clusterSources(),
      this.$coreAPI.namespacedSources()
    ])

    this.clusterPages = clusterSources.map(s => ({ title: s, path: `/cluster-policy-reports/${s}` }))
    this.namespacedPages = namespacedSources.map(s => ({ title: s, path: `/policy-reports/${s}` }))
  },
  computed: {
    ...mapGetters(['refreshInterval', 'currentCluster']),
    showLogs (): boolean {
      if (!this.views.logs) {
        return false
      }

      if (!this.currentCluster) {
        return true
      }

      return this.currentCluster.id === ''
    }
  },
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        if (!refreshInterval) {
          this.interval = null
          this.$fetch()
          return
        }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    },
    currentCluster: '$fetch'
  },
  methods: {
    showPlugin (plugin: Plugin): boolean {
      if (!this.plugins || !this.plugins.includes(plugin)) {
        return false
      }

      if (this.currentCluster) {
        return this.currentCluster[plugin] as boolean
      }

      return true
    }
  }
})
</script>
