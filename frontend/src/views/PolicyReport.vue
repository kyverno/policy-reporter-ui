<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
        <v-col>
            <v-toolbar elevation="1">
                <div style="width: 550px;" class="mr-2">
                  <policy-autocomplete v-model="policies" :policies="availablePolicies" />
                </div>
                <v-spacer />
                <div style="width: 450px;">
                  <namespace-autocomplete v-model="namespaces" />
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

const flatResults = (policies: string[], namespaces: string[], reports: GlobalPolicyReportMap) => policies.reduce<Result[]>((acc, policy) => {
  if (!reports[policy]) {
    return acc;
  }

  if (namespaces.length > 0) {
    return [...acc, ...reports[policy].results.filter((result) => namespaces.includes(result.resource.namespace as string))];
  }

  return [...acc, ...reports[policy].results];
}, []);

type Data = {
  heights: { [key in Status]: number };
  policies: string[];
  namespaces: string[];
}

type Methods = {
  updateHeight(status: string, height: number): void;
  statusText(status: string): string;
}

type Computed = {
  globalPolicyMap: GlobalPolicyReportMap;
  availablePolicies: string[];
  results: Result[];
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
    results(): Result[] {
      return flatResults(this.policies, this.namespaces, this.globalPolicyMap);
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
