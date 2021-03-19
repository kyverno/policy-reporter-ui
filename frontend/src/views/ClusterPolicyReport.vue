<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
        <v-col>
            <v-toolbar elevation="1">
                <div style="width: 550px;">
                  <policy-autocomplete v-model="policies" :policies="availablePolicies" />
                </div>
            </v-toolbar>
        </v-col>
    </v-row>
    <v-row v-if="policies.length === 0">
      <v-col>
        <v-card>
          <v-card-text>
            <v-alert type="info" outlined class="ma-0">
              Select one ore more Policies to get related results.
            </v-alert>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6" md="3">
        <cluster-policy-status color="green" :count="passingResults.length" :statusText="statusText('pass')" />
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <cluster-policy-status color="warning" :count="warningResults.length" :statusText="statusText('warn')" />
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <cluster-policy-status color="error" :count="failingResults.length" :statusText="statusText('fail')" />
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <cluster-policy-status color="error" :count="errorResults.length" :statusText="statusText('error')" />
      </v-col>
    </v-row>
    <v-row v-if="errorResults.length > 0">
      <v-col cols="12">
        <cluster-policy-table :results="errorResults" title="Error ClusterPolicy Results" />
      </v-col>
    </v-row>
    <v-row v-if="failingResults.length > 0">
      <v-col cols="12">
        <cluster-policy-table :results="failingResults" title="Failing ClusterPolicy Results" />
      </v-col>
    </v-row>
    <v-row v-if="warningResults.length > 0">
      <v-col cols="12">
        <cluster-policy-table :results="warningResults" title="Warning ClusterPolicy Results" />
      </v-col>
    </v-row>
    <v-row v-if="passingResults.length > 0">
      <v-col cols="12">
        <cluster-policy-table :results="passingResults" title="Passing ClusterPolicy Results" />
      </v-col>
    </v-row>
    <v-row v-if="skippedResults.length > 0">
      <v-col cols="12">
        <cluster-policy-table :results="skippedResults" title="Skipping ClusterPolicy Results" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import { mapState } from 'vuex';
import { ClusterPolicyReport, Result, Status } from '@/models';
import PolicyAutocomplete from '@/components/PolicyAutocomplete.vue';
import ClusterPolicyTable from '@/components/ClusterPolicyTable.vue';
import ClusterPolicyStatus from '@/components/ClusterPolicyStatus.vue';
import { flatPolicies } from '@/mapper';

const flatResults = (policies: string[], reports: Array<ClusterPolicyReport>) => reports.reduce<Result[]>((acc, item) => {
  item.results.forEach((result: Result) => {
    if (!policies.includes(result.policy)) {
      return;
    }

    acc.push(result);
  });

  return acc;
}, []);

type Data = { minHeight: number }
type Methods = { updateHeight(height: number): void; statusText(status: string): string }
type Computed = {
  clusterReports: ClusterPolicyReport[];
  availablePolicies: string[];
  results: Result[];
  skippedResults: Result[];
  passingResults: Result[];
  warningResults: Result[];
  failingResults: Result[];
  errorResults: Result[];
  policies: string[];
}

export default Vue.extend<Data, Methods, Computed, {}>({
  components: {
    PolicyAutocomplete,
    ClusterPolicyTable,
    ClusterPolicyStatus,
  },
  name: 'PolicyReport',
  data: () => ({
    minHeight: 0,
  }),
  watch: {
    policies() {
      this.minHeight = 0;
    },
  },
  computed: {
    ...mapState(['clusterReports']),
    availablePolicies() {
      const policies = flatPolicies(this.clusterReports);

      policies.sort();

      return policies;
    },
    results(): Result[] {
      return flatResults(this.policies, this.clusterReports);
    },
    skippedResults() {
      return this.results.filter((r) => r.status === Status.SKIP);
    },
    passingResults() {
      return this.results.filter((r) => r.status === Status.PASS);
    },
    warningResults() {
      return this.results.filter((r) => r.status === Status.WARN);
    },
    failingResults() {
      return this.results.filter((r) => r.status === Status.FAIL);
    },
    errorResults() {
      return this.results.filter((r) => r.status === Status.ERROR);
    },
    policies: {
      get(): string[] {
        if (Array.isArray(this.$route.query.policies)) {
          return this.$route.query.policies as string[];
        }

        if (!this.$route.query.policies) {
          return [];
        }

        return [this.$route.query.policies as string];
      },
      set(policies: string[]) {
        this.$router.push({ name: this.$route.name as string, query: { policies } });
      },
    },
  },
  methods: {
    updateHeight(height: number) {
      if (height < this.minHeight) return;

      this.minHeight = height;
    },
    statusText(status: string): string {
      switch (status) {
        case Status.SKIP:
          return 'Skipped';
        case Status.PASS:
          return 'Passing';
        case Status.WARN:
          return 'Warning';
        case Status.FAIL:
          return 'Failing';
        case Status.ERROR:
          return 'Errored';
        default:
          return 'Unknown';
      }
    },
  },
});
</script>
