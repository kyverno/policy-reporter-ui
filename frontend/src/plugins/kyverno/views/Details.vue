<template>
  <v-container fluid class="py-8 px-6" v-if="policy">
    <v-row>
      <v-col>
      <v-card width="100%">
        <v-card-title>
          {{ policy.kind }} "{{ policy.name }}"
          <v-spacer />
          <v-chip class="mb-1" color="green" text-color="white" v-if="policy.background">Background scan enabled</v-chip>
        </v-card-title>
        <template  v-if="policy.description">
          <v-divider />
          <v-card-text>
            {{ policy.description }}
          </v-card-text>
        </template>
      </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="3">
        <policy-details :policy="policy" />
      </v-col>
      <v-col cols="3">
        <v-card style="height: 100%">
          <v-toolbar flat>
            <v-toolbar-title>Configured Autogen Controller</v-toolbar-title>
          </v-toolbar>
          <v-list>
            <template v-for="controller in policy.autogenControllers" >
            <v-divider :key="controller + '_divider'" />
            <v-list-item :key="controller">
              <v-list-item-content>
                <v-list-item-title>
                  {{ controller }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            </template>
          </v-list>
        </v-card>
      </v-col>
      <v-col cols="3">
        <v-card min-height="300" style="height: 100%">
          <v-card-title class="pb-0">
            Passed Results
          </v-card-title>
          <v-card-text style="height: calc(100% - 48px); font-size: 9rem!important;"
                      :class="['text-center text-h1 d-flex justify-center align-center', `success--text`]"
          >
            {{ results.pass.length }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="3">
        <v-card min-height="300" style="height: 100%">
          <v-card-title class="pb-0">
            Failed Results
          </v-card-title>
          <v-card-text style="height: calc(100% - 48px); font-size: 9rem!important;"
                      :class="['text-center text-h1 d-flex justify-center align-center', `error--text`]"
          >
            {{ results.fail.length }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <rules-table :policy="policy" :results="results" />
      </v-col>
    </v-row>
    <v-row v-if="validations">
      <v-col cols="12">
        <policy-table :results="results.fail" title="Failed Results" />
      </v-col>
    </v-row>
    <v-row v-if="validations">
      <v-col cols="12">
        <policy-table :results="results.pass" title="Passed Results" />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-toolbar flat>
            <v-toolbar-title>
              YAML File
            </v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <highlightjs language='yaml' :code="policy.content" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-btn @click="$router.push('/kyverno')" fixed top right style="z-index: 999;" depressed dark color="grey darken-1">
      Back
    </v-btn>
  </v-container>
</template>

<script lang="ts">
import PolicyTable from '@/components/PolicyTable.vue';
import { GlobalPolicyReportMap, Status } from '@/models';
import Vue from 'vue';
import { mapGetters, mapState } from 'vuex';
import PolicyDetails from '../components/PolicyDetails.vue';
import RulesTable from '../components/RulesTable.vue';
import { Policy, ResultMap, RuleType } from '../models';
import { NAMESPACE, FETCH_POLICIES } from '../store';

type Data = {}
type Methods = {}
type Props = {
  uid: string;
}
type Computed = {
    policy?: Policy;
    policies: Policy[];
    globalPolicyMap: GlobalPolicyReportMap;
    results: ResultMap;
    validations: boolean;
}

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { PolicyDetails, RulesTable, PolicyTable },
  name: 'Details',
  props: {
    uid: { type: String, required: true },
  },
  data: () => ({}),
  computed: {
    ...mapGetters(NAMESPACE, ['policies']),
    ...mapState(['globalPolicyMap']),
    policy() {
      return this.policies.find((p) => p.uid === this.uid);
    },
    validations() {
      if (!this.policy) return false;

      return this.policy.rules.some((r) => r.type === RuleType.VALIDATION);
    },
    results(): ResultMap {
      const defaults: ResultMap = { fail: [], pass: [] };

      if (!this.policy) return defaults;

      const results = this.globalPolicyMap[this.policy.name];
      if (!results) return defaults;

      results.results.forEach((result) => {
        if (result.status === Status.FAIL) {
          defaults.fail.push(result);
        }
        if (result.status === Status.PASS) {
          defaults.pass.push(result);
        }
      });

      return defaults;
    },
  },
  created() {
    this.$store.dispatch(`${NAMESPACE}/${FETCH_POLICIES}`);
  },
});
</script>

<style scoped>
>>> code {
  padding: 16px!important;
}
</style>
