<template>
  <v-navigation-drawer :value="value" app clipped width="305" @input="$emit('input', $event)">
    <v-list>
      <app-navigation-item icon="mdi-view-dashboard" route="/" title="Dashboard" />
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

      <app-navigation-item icon="mdi-console" route="/logs" title="Logs" />
      <app-navigation-item v-if="plugins && plugins.includes('kyverno')" route="/kyverno-plugin" title="Kyverno Policies" exact>
        <template #icon>
          <v-list-item-icon>
            <lazy-kyverno-icon style="height: 24px; width: 24px;" />
          </v-list-item-icon>
        </template>
      </app-navigation-item>

      <app-navigation-item
        v-if="plugins && plugins.includes('kyverno')"
        route="/kyverno-plugin/verify-image-rules"
        title="Kyverno ImageVerify Rules"
        icon="mdi-shield-check"
      />
    </v-list>
  </v-navigation-drawer>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import AppNavigationItem from './AppNavigationItem.vue'

type Page = { title: string; path: string }

type Data = {
  interval: any,
  clusterPages: Page[];
  namespacedPages: Page[];
}

type Props = {
  value: boolean;
  plugins: string[];
}

type Computed = {
  refreshInterval: number;
}

type Methdos = {}

export default Vue.extend<Data, Methdos, Computed, Props>({
  components: { AppNavigationItem },
  props: {
    value: { type: Boolean, default: false },
    plugins: { type: Array, default: () => [] }
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
  computed: mapGetters(['refreshInterval']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  }
})
</script>
