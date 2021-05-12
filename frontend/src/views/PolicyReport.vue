<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
        <v-col>
            <v-card elevation="1">
              <v-container fluid>
                <v-row>
                  <v-col cols="4" class="d-inline-block">
                    <policy-autocomplete v-model="policies" :policies="availablePolicies" />
                  </v-col>
                  <v-col cols="4">
                    <kind-autocomplete v-model="kinds" :kinds="availableKinds" />
                  </v-col>
                  <v-col cols="4" class="d-inline-block">
                    <namespace-autocomplete v-model="namespaces" />
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="4">
                    <category-autocomplete v-model="categories" :categories="availableCategories" />
                  </v-col>
                  <v-col cols="4">
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
      <policy-status-per-namespace @height-change="updateHeight('fail', $event)"
                                  :minHeight="minHeight"
                                  :results="failingResults"
                                  :statusText="statusText('fail')"
      />
      <policy-status-per-namespace @height-change="updateHeight('pass', $event)"
                                  :minHeight="minHeight"
                                  :results="passingResults"
                                  :statusText="statusText('pass')"
      />
      <policy-status-per-namespace @height-change="updateHeight('error', $event)"
                                  :minHeight="minHeight"
                                  :results="errorResults"
                                  :statusText="statusText('error')"
                                  optional
      />
      <policy-status-per-namespace @height-change="updateHeight('warn', $event)"
                                  :minHeight="minHeight"
                                  :results="warningResults"
                                  :statusText="statusText('warn')"
                                  optional
      />
      <policy-status-per-namespace @height-change="updateHeight('skip', $event)"
                                  :minHeight="minHeight"
                                  :results="skippedResults"
                                  :statusText="statusText('skip')"
                                  optional
      />
    </v-row>
    <v-row v-if="errorResults.length > 0">
      <v-col cols="12">
        <policy-table :results="errorResults" title="Error Policy Results" />
      </v-col>
    </v-row>
    <v-row v-if="failingResults.length > 0">
      <v-col cols="12">
        <policy-table :results="failingResults" title="Failing Policy Results" />
      </v-col>
    </v-row>
    <v-row v-if="warningResults.length > 0">
      <v-col cols="12">
        <policy-table :results="warningResults" title="Warning Policy Results" />
      </v-col>
    </v-row>
    <v-row v-if="passingResults.length > 0">
      <v-col cols="12">
        <policy-table :results="passingResults" title="Passing Policy Results" />
      </v-col>
    </v-row>
    <v-row v-if="skippedResults.length > 0">
      <v-col cols="12">
        <policy-table :results="skippedResults" title="Skipping Policy Results" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import { mapState } from 'vuex';
import { GlobalPolicyReportMap, Result, Status } from '@/models';
import PolicyStatusPerNamespace from '@/components/PolicyStatusPerNamespace.vue';
import PolicyTable from '@/components/PolicyTable.vue';
import PolicyAutocomplete from '@/components/PolicyAutocomplete.vue';
import NamespaceAutocomplete from '@/components/NamespaceAutocomplete.vue';
import CategoryAutocomplete from '@/components/CategoryAutocomplete.vue';
import SeverityAutocomplete from '@/components/SeverityAutocomplete.vue';
import KindAutocomplete from '@/components/KindAutocomplete.vue';

const flatResults = (
  policies: string[],
  reports: GlobalPolicyReportMap,
) => policies.reduce<Result[]>((acc, policy) => {
  if (!reports[policy]) {
    return acc;
  }

  return [...acc, ...reports[policy].results];
}, []);

type Data = {
  heights: { [key in Status]: number };
  policies: string[];
  namespaces: string[];
  categories: string[];
  severities: string[];
  kinds: string[];
}

type Methods = {
  updateHeight(status: string, height: number): void;
  statusText(status: string): string;
}

type Computed = {
  globalPolicyMap: GlobalPolicyReportMap;
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
  minHeight: number;
}

export default Vue.extend<Data, Methods, Computed, {}>({
  components: {
    PolicyStatusPerNamespace,
    PolicyTable,
    PolicyAutocomplete,
    NamespaceAutocomplete,
    CategoryAutocomplete,
    SeverityAutocomplete,
    KindAutocomplete,
  },
  name: 'PolicyReport',
  data: () => ({
    heights: {
      [Status.SKIP]: 0,
      [Status.PASS]: 0,
      [Status.WARN]: 0,
      [Status.FAIL]: 0,
      [Status.ERROR]: 0,
    },
    policies: [],
    namespaces: [],
    categories: [],
    severities: [],
    kinds: [],
  }),
  computed: {
    ...mapState(['globalPolicyMap']),
    minHeight() {
      return Object.values(this.heights).reduce<number>((acc, height) => {
        if (acc > height) return acc;

        return height;
      }, 0);
    },
    availablePolicies(): string[] {
      return Object.keys(this.globalPolicyMap);
    },
    availableCategories(): string[] {
      return Object.keys(this.results.reduce<{ [category: string]: boolean }>((c, r) => {
        if (!r.category) return c;

        return { ...c, [r.category]: true };
      }, {}));
    },
    availableKinds(): string[] {
      return Object.keys(this.results.reduce<{ [kind: string]: boolean }>((c, r) => {
        if (!r.resource) {
          return c;
        }

        return { ...c, [r.resource.kind]: true };
      }, {}));
    },
    results(): Result[] {
      return flatResults(this.policies, this.globalPolicyMap);
    },
    filteredResults(): Result[] {
      return this.results.filter((result) => {
        if (this.namespaces.length > 0 && result.resource && !this.namespaces.includes(result.resource.namespace as string)) return false;

        if (this.kinds.length > 0 && result.resource && !this.kinds.includes(result.resource.kind)) return false;

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
    updateHeight(status: string, height: number) {
      this.heights = { ...this.heights, [status as Status]: height };
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
