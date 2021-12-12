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
import Vue from 'vue'
import { DataTableHeader } from 'vuetify'
import { Policy, Rule } from '../types'
import { Status } from '~/policy-reporter-plugins/core/types'

type Item = Rule & { passed: number; failed: number };

export default Vue.extend<{ items: Item[] }, {}, { headers: DataTableHeader[] }, { policy: Policy; }>({
  props: {
    policy: { type: Object, required: true }
  },
  data: () => ({ items: [] }),
  async fetch () {
    this.items = await Promise.all(this.policy.rules.map<Promise<Item>>(async (rule) => {
      let failed: number = 0
      let passed: number = 0

      const counts = await this.$coreAPI.ruleStatusCount(this.policy.name, rule.name)
      counts.forEach((result) => {
        if (result.status === Status.PASS) {
          passed += result.count
        }

        if (result.status === Status.FAIL) {
          failed += result.count
        }
      })

      return { ...rule, passed, failed }
    }))
  },
  computed: {
    headers (): DataTableHeader[] {
      return [
        { text: 'Name', value: 'name' },
        { text: 'Type', value: 'type' },
        { text: 'Message', value: 'type' },
        { text: 'Passed', value: '' },
        { text: 'Failed', value: '' }
      ]
    }
  }
})
</script>
