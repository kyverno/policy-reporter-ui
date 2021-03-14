<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
      <v-col cols="12" md="8">
        <v-card>
          <failing-per-namespace :reports="namespacePolicyMap" />
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <failing-cluster-policies :reports="clusterReports" />
      </v-col>
    </v-row>
    <v-row v-if="results.clusterReportResults.length > 0">
      <v-col cols="12">
        <cluster-policy-table :results="results.clusterReportResults" title="Failing ClusterPolicy Results" />
      </v-col>
    </v-row>
    <v-row v-if="results.reportResults.length > 0">
      <v-col cols="12">
        <policy-table :results="results.reportResults" title="Failing Policy Results" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import { mapState } from 'vuex';
import FailingPerNamespace from '@/components/FailingPerNamespace.vue';
import FailingClusterPolicies from '@/components/FailingClusterPolicies.vue';
import {
  ClusterPolicyReport, NamespacePolicyReport, Result, Status,
} from '@/models';
import ClusterPolicyTable from '@/components/ClusterPolicyTable.vue';
import PolicyTable from '@/components/PolicyTable.vue';

const convertReports = (reports: Array<NamespacePolicyReport|ClusterPolicyReport>) => reports.reduce<Result[]>((acc, item) => {
  item.results.forEach((result: Result) => {
    if (![Status.FAIL, Status.ERROR].includes(result.status)) {
      return;
    }

    acc.push(result);
  });

  return acc;
}, []);

export default Vue.extend({
  components: {
    FailingPerNamespace, FailingClusterPolicies, ClusterPolicyTable, PolicyTable,
  },
  name: 'Dashboard',
  data: () => ({}),
  computed: {
    ...mapState(['namespacePolicyMap', 'clusterReports']),
    results(): { [type: string]: Result[] } {
      return {
        reportResults: convertReports(Object.values(this.namespacePolicyMap)),
        clusterReportResults: convertReports(this.clusterReports),
      };
    },
  },
});
</script>
