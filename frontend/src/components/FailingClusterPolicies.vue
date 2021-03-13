<template>
  <v-card full-height style="height: 100%" min-height="300">
    <v-card-title class="pb-0">
      Failing Cluster Policies
    </v-card-title>
    <v-card-text style="height: calc(100% - 48px); font-size: 11rem!important;" class="text-center text-h1 d-flex justify-center align-center">
      {{ failing }}
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue';
import { ClusterPolicyReport } from '@/models';

export default Vue.extend<{}, {}, { failing: number }, { reports: ClusterPolicyReport[] }>({
  name: 'FailingClusterPolicies',
  props: {
    reports: { required: true, type: Array },
  },
  computed: {
    failing() {
      return (this.reports as ClusterPolicyReport[]).reduce<number>((acc, item) => acc + item.summary.fail, 0);
    },
  },
});
</script>
