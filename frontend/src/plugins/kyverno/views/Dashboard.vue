<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
      <v-col cols="12" md="6" v-if="policyGroups">
        <category-chart :policyGroups="policyGroups" />
      </v-col>
      <v-col cols="12" md="6">
        <rule-type-chart :policies="policies" />
      </v-col>
    </v-row>
    <template  v-for="(group, category) in policyGroups">
      <v-row :key="category" v-if="group.length > 0">
        <v-col cols="12">
          <policy-table :policies="group" :title="category" />
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import { mapGetters } from 'vuex';
import { PolicyGroups, Policy } from '../models';
import PolicyTable from '../components/PolicyTable.vue';
import { NAMESPACE } from '../store';
import RuleTypeChart from '../components/RuleTypeChart.vue';
import CategoryChart from '../components/CategoryChart.vue';

type Data = {}
type Methods = {}
type Props = {}
type Computed = {
    policyGroups: PolicyGroups;
    policies: Policy[];
}

export default Vue.extend<Data, Methods, Computed, Props>({
  components: {
    PolicyTable,
    RuleTypeChart,
    CategoryChart,
  },
  name: 'PolicyReport',
  data: () => ({}),
  computed: {
    ...mapGetters(NAMESPACE, ['policyGroups', 'policies']),
  },
});
</script>
