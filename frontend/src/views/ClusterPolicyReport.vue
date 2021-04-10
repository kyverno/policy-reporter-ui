<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
        <v-col>
            <v-card elevation="1">
              <v-container fluid>
                <v-row>
                  <v-col cols="6" class="d-inline-block">
                    <policy-autocomplete v-model="policies" :policies="availablePolicies" />
                  </v-col>
                  <v-col cols="6">
                    <kind-autocomplete v-model="kinds" :kinds="availableKinds" />
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="6">
                    <category-autocomplete v-model="categories" :categories="availableCategories" />
                  </v-col>
                  <v-col cols="6">
                    <severity-autocomplete v-model="severities" />
                  </v-col>
                </v-row>
              </v-container>
            </v-card>
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
import CategoryAutocomplete from '@/components/CategoryAutocomplete.vue';
import SeverityAutocomplete from '@/components/SeverityAutocomplete.vue';
import KindAutocomplete from '@/components/KindAutocomplete.vue';

const flatResults = (policies: string[], reports: Array<ClusterPolicyReport>) => reports.reduce<Result[]>((acc, item) => {
  item.results.forEach((result: Result) => {
    if (!policies.includes(result.policy)) {
      return;
    }

    acc.push(result);
  });

  return acc;
}, []);

type Data = {
  minHeight: number;
  policies: string[];
  categories: string[];
  severities: string[];
  kinds: string[];
}
type Methods = { updateHeight(height: number): void; statusText(status: string): string }
type Computed = {
  clusterReports: ClusterPolicyReport[];
  availablePolicies: string[];
  availableCategories: string[];
  availableKinds: string[];
  results: Result[];
  filteredResults: Result[];
  skippedResults: Result[];
  passingResults: Result[];
  warningResults: Result[];
  failingResults: Result[];
  errorResults: Result[];
}

export default Vue.extend<Data, Methods, Computed, {}>({
  components: {
    PolicyAutocomplete,
    ClusterPolicyTable,
    ClusterPolicyStatus,
    CategoryAutocomplete,
    SeverityAutocomplete,
    KindAutocomplete,
  },
  name: 'PolicyReport',
  data: () => ({
    minHeight: 0,
    policies: [],
    categories: [],
    severities: [],
    kinds: [],
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
    availableCategories(): string[] {
      return Object.keys(this.results.reduce<{ [category: string]: boolean }>((c, r) => {
        if (!r.category) return c;

        return { ...c, [r.category]: true };
      }, {}));
    },
    availableKinds(): string[] {
      return Object.keys(this.results.reduce<{ [kind: string]: boolean }>((c, r) => ({ ...c, [r.resource.kind]: true }), {}));
    },
    filteredResults(): Result[] {
      return this.results.filter((result) => {
        if (this.kinds.length > 0 && !this.kinds.includes(result.resource.kind)) return false;

        if (this.categories.length > 0 && !this.categories.includes(result.category || '')) return false;

        if (this.severities.length > 0 && !this.severities.includes(result.severity || '')) return false;

        return true;
      });
    },
    skippedResults() {
      return this.filteredResults.filter((r) => r.status === Status.SKIP);
    },
    passingResults() {
      return this.filteredResults.filter((r) => r.status === Status.PASS);
    },
    warningResults() {
      return this.filteredResults.filter((r) => r.status === Status.WARN);
    },
    failingResults() {
      return this.filteredResults.filter((r) => r.status === Status.FAIL);
    },
    errorResults() {
      return this.filteredResults.filter((r) => r.status === Status.ERROR);
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
