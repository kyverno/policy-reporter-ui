<template>
    <v-autocomplete dense
                    multiple
                    :items="items"
                    outlined
                    hide-details
                    label="Policies"
                    :value="value"
                    v-bind="$attrs"
                    v-on="$listeners"
    >
    <template v-slot:selection="{ item, index }">
        <v-chip small v-if="index <= 1" label outlined>
          <span>{{ item }}</span>
        </v-chip>
        <span
          v-if="index === 2"
          class="grey--text caption"
        >
          (+{{ value.length - 2 }} others)
        </span>
      </template>
    </v-autocomplete>
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
