<template>
  <v-card>
    <v-toolbar flat>
      <v-toolbar-title>Rules</v-toolbar-title>
    </v-toolbar>
    <v-divider />
    <v-data-table
        :items="items"
        :headers="headers"
        :items-per-page="10"
        :sort-by="['name']"
        item-key="name"
    >
    <template #item="{ item }">
        <tr>
            <td>{{ item.name }}</td>
            <td>{{ item.type }}</td>
            <td>{{ item.message }}</td>
            <td>{{ item.passed }}</td>
            <td>{{ item.failed }}</td>
        </tr>
    </template>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue';
import { DataTableHeader } from 'vuetify';
import { Policy, ResultMap, Rule } from '../models';

type Item = Rule & { passed: number; failed: number };

export default Vue.extend<{}, {}, { headers: DataTableHeader[]; items: Item[] }, { policy: Policy; results: ResultMap }>({
  props: {
    policy: { type: Object, required: true },
    results: { type: Object, required: true },
  },
  computed: {
    items(): Item[] {
      return this.policy.rules.map<Item>((rule) => {
        let failed = 0;
        let passed = 0;

        this.results.pass.forEach((result) => {
          if (result.rule !== rule.name) return;

          passed += 1;
        });
        this.results.fail.forEach((result) => {
          if (result.rule !== rule.name) return;

          failed += 1;
        });

        return { ...rule, failed, passed };
      });
    },
    headers(): DataTableHeader[] {
      return [
        { text: 'Name', value: 'name' },
        { text: 'Type', value: 'type' },
        { text: 'Message', value: 'type' },
        { text: 'Passed', value: '' },
        { text: 'Failed', value: '' },
      ];
    },
  },
});
</script>
