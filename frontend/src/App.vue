<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" app clipped>
      <v-list>
        <template v-for="[icon, text, route, plugin, customIcon] in links">
        <v-list-item
          :key="route"
          :to="route"
          v-if="!plugin || plugins.includes(plugin)"
        >
          <v-list-item-icon>
            <v-icon v-if="icon">{{ icon }}</v-icon>
            <component style="height: 24px; width: 24px;" v-else :is="customIcon"></component>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ text }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app class="elevation-1" clipped-left>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>

      <v-toolbar-title>Policy Reporter</v-toolbar-title>

      <v-spacer />

      <template v-if="$route.name !== 'policy-details'">
        <v-chip v-for="target in targets" class="ml-3" outlined :key="target.name" label :title="`minimum priority: ${target.minimumPriority}`">
          <v-avatar left>
            <v-icon :color="target.minimumPriority | mapPriority">mdi-target-variant</v-icon>
          </v-avatar>
          {{ target.name }}
        </v-chip>
      </template>
    </v-app-bar>

    <v-main class="">
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import { mapState } from 'vuex';
import {
  FETCH_TARGETS, FETCH_REPORTS, FETCH_CLUSTER_REPORTS, FETCH_PLUGINS,
} from '@/store';

const KyvernoIcon = () => import('@/plugins/kyverno/components/KyvernoIcon.vue');

export default Vue.extend({
  components: { KyvernoIcon },
  data: () => ({
    drawer: null,
    interval: 0,
    links: [
      ['mdi-view-dashboard', 'Dashboard', '/'],
      ['mdi-file-chart', 'Policy Reports', '/policy-reports'],
      ['mdi-file-chart', 'ClusterPolicy Reports', '/cluster-policy-reports'],
      ['mdi-console', 'Logs', '/logs'],
      [null, 'Kyverno Policies', '/kyverno', 'kyverno', 'kyverno-icon'],
    ],
  }),
  computed: mapState(['targets', 'plugins']),
  created() {
    this.$store.dispatch(FETCH_PLUGINS);
    this.$store.dispatch(FETCH_TARGETS);
    this.$store.dispatch(FETCH_REPORTS);
    this.$store.dispatch(FETCH_CLUSTER_REPORTS);

    // @ts-ignore
    this.interval = setInterval(() => {
      this.$store.dispatch(FETCH_REPORTS);
      this.$store.dispatch(FETCH_CLUSTER_REPORTS);
    }, 5000);
  },
  destroyed() {
    clearInterval(this.interval);
  },
});
</script>

<style>
.apexcharts-svg {
  background-color: transparent!important;
}
</style>
