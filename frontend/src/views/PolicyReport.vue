<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
        <v-col>
            <v-toolbar elevation="1">
                <div style="max-width: 600px;">
                  <policy-autocomplete v-model="policies" :reports="reports" />
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
      <policy-status-per-namespace @height-change="updateHeight($event)"
                                  :minHeight="minHeight"
                                  :results="passingResults"
                                  :statusText="statusText('pass')"
      />
      <policy-status-per-namespace @height-change="updateHeight($event)"
                                  :minHeight="minHeight"
                                  :results="failingResults"
                                  :statusText="statusText('fail')"
      />
      <policy-status-per-namespace @height-change="updateHeight($event)"
                                  :minHeight="minHeight"
                                  :results="warningResults"
                                  :statusText="statusText('warn')"
                                  optional
      />
      <policy-status-per-namespace @height-change="updateHeight($event)"
                                  :minHeight="minHeight"
                                  :results="errorResults"
                                  :statusText="statusText('error')"
                                  optional
      />
      <policy-status-per-namespace @height-change="updateHeight($event)"
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
import { PolicyReport, Result, Status } from '@/models';
import PolicyStatusPerNamespace from '@/components/PolicyStatusPerNamespace.vue';
import PolicyTable from '@/components/PolicyTable.vue';
import PolicyAutocomplete from '@/components/PolicyAutocomplete.vue';

const flatResults = (policies: string[], reports: Array<PolicyReport>) => reports.reduce<Result[]>((acc, item) => {
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
  reports: PolicyReport[];
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
    PolicyStatusPerNamespace,
    PolicyTable,
    PolicyAutocomplete,
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
    ...mapState(['reports']),
    results(): Result[] {
      return flatResults(this.policies, this.reports);
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
