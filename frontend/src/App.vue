<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" app clipped>
      <v-list>
        <v-list-item
          v-for="[icon, text, route] in links"
          :key="route"
          :to="route"
        >
          <v-list-item-icon>
            <v-icon>{{ icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ text }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app class="elevation-1" clipped-left>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>

      <v-toolbar-title>Policy Reporter</v-toolbar-title>

      <v-spacer />

      <v-chip v-for="target in targets" class="ml-3" outlined :key="target.name" label :title="`minimum priority: ${target.minimumPriority}`">
        <v-avatar left>
          <v-icon :color="target.minimumPriority | mapPriority">mdi-target-variant</v-icon>
        </v-avatar>
        {{ target.name }}
      </v-chip>
    </v-app-bar>

    <v-main class="grey lighten-5">
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue';
import { mapState } from 'vuex';
import {
  FETCH_TARGETS, FETCH_REPORTS, FETCH_CLUSTER_REPORTS,
} from '@/store';

export default Vue.extend({
  data: () => ({
    cards: ['Today', 'Yesterday'],
    drawer: null,
    interval: 0,
    links: [
      ['mdi-view-dashboard', 'Dashboard', '/'],
      ['mdi-file-chart', 'Policy Reports', '/policy-reports'],
      ['mdi-file-chart', 'ClusterPolicy Reports', '/cluster-policy-reports'],
      ['mdi-console', 'Logs', '/logs'],
    ],
  }),
  computed: mapState(['targets']),
  created() {
    this.$store.dispatch(FETCH_TARGETS);
    this.$store.dispatch(FETCH_REPORTS);
    this.$store.dispatch(FETCH_CLUSTER_REPORTS);

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
