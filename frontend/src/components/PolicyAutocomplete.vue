<template>
    <v-autocomplete dense multiple :items="items" outlined hide-details label="Policies" :value="value" v-bind="$attrs" v-on="$listeners" />
</template>

<script lang="ts">
import { PolicyReport, ClusterPolicyReport, Result } from '@/models';
import Vue from 'vue';

const flatPolicies = (reports: Array<PolicyReport|ClusterPolicyReport>) => reports.reduce<string[]>((acc, item) => {
  item.results.forEach((result: Result) => {
    if (acc.includes(result.policy)) {
      return;
    }

    acc.push(result.policy);
  });

  return acc;
}, []);

export default Vue.extend<{}, {}, { items: string[] }, { reports: Array<PolicyReport|ClusterPolicyReport>; value: string[] }>({
  props: {
    reports: { type: Array, default: () => [] },
    value: { type: Array, default: () => [] },
  },
  computed: {
    items() {
      const policies = flatPolicies(this.reports);

      policies.sort();

      return policies;
    },
  },
  created() {
    if (!this.items.length || this.value.length > 0) return;

    this.$emit('input', [this.items[0]]);
  },
});
</script>
