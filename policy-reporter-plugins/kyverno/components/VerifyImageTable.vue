<template>
  <v-card>
    <v-card-title>
      {{ title }}
      <v-spacer />
      <v-btn icon @click="open = !open">
        <v-icon>{{ open ? 'mdi-chevron-up' : 'mdi-chevron-down' }}</v-icon>
      </v-btn>
    </v-card-title>
    <v-expand-transition>
      <div v-show="open">
        <v-divider />
        <v-card-title>
          <v-spacer />
          <div style="width: 450px;">
            <v-text-field
              v-model="search"
              append-icon="mdi-magnify"
              label="Search"
              outlined
              dense
              hide-details
              clearable
            />
          </div>
        </v-card-title>
        <v-divider />
        <v-data-table
          :items="items"
          :headers="tableHeaders"
          :items-per-page="10"
          :search="search"
          :sort-by="['namespace', 'name', 'repository', 'image', 'uid']"
          item-key="uid"
        >
          <template #item="{ item }">
            <tr>
              <td v-if="showNamespace">
                {{ item.namespace }}
              </td>
              <td>{{ item.name }}</td>
              <td>{{ item.rule }}</td>
              <td>{{ item.repository }}</td>
              <td>{{ item.image }}</td>
            </tr>
          </template>
        </v-data-table>
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { DataTableHeader } from 'vuetify'
import { Policy, VerifyImageRule } from '../types'

type Item = Omit<VerifyImageRule, 'policy'> & { name: string; namespace?: string; uid: string }

type Data = { open: boolean; search: string; }
type Computed = { tableHeaders: DataTableHeader[]; showNamespace: boolean; items: Item[] }
type Props = { title: string; rules: VerifyImageRule[] }

export default Vue.extend<Data, {}, Computed, Props>({
  props: {
    title: { type: String, required: true },
    rules: { type: Array, required: true }
  },
  data: () => ({ open: true, search: '', expanded: [] }),
  computed: {
    items () {
      return this.rules.map<Item>(({ policy, ...r }) => ({ ...r, ...policy }))
    },
    showNamespace (): boolean {
      return this.items.some(p => !!p.namespace)
    },
    tableHeaders (): DataTableHeader[] {
      return [
        ...(this.showNamespace ? [{ text: 'Namespace', value: 'namespace' }] : []),
        { text: 'Policy', value: 'name' },
        { text: 'Rule', value: 'rule' },
        { text: 'Cosign Repository', value: 'repository' },
        { text: 'Image', value: 'image' }
      ]
    }
  },
  methods: {
    types (policy: Policy): string {
      return policy.rules
        .map(r => r.type)
        .filter((type, index, self) => self.indexOf(type) === index).join(', ')
    }
  }
})
</script>

<style scoped>
>>> code {
  padding: 16px !important;
}
</style>
