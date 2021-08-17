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
              :headers="headers"
              :items-per-page="10"
              :search="search"
              :sort-by="['namespace', 'name']"
              :expanded.sync="expanded"
              item-key="uid"
              >
            <template #item="{ item, expand, isExpanded }">
                <tr @click.stop="expand(!isExpanded)" style="cursor: pointer">
                  <td>{{ item.kind }}</td>
                  <td v-if="showNamespace">{{ item.namespace }}</td>
                  <td>{{ item.name }}</td>
                  <td>{{ item.rules.length }}</td>
                  <td>{{ item.type }}</td>
                  <td>{{ item.background }}</td>
                  <td>{{ item.validationFailureAction }}</td>
                  <td>
                    <severity-chip :severity="item.severity" label @click.stop="search = item.severity" v-if="item.severity" />
                  </td>
                  <td>
                        <v-btn dark @click.stop="$router.push(`/kyverno-plugin/${item.uid}`)" class="mr-2" depressed small>
                          Details
                        </v-btn>
                        <yaml-dialog :policy="item" />
                  </td>
                </tr>
            </template>
            <template #expanded-item="{ headers, item }">
              <tr class="grey lighten-4">
                <td :colspan="headers.length" class="py-3">
                  {{ item.description }}
                </td>
              </tr>
            </template>
            </v-data-table>
        </div>

        </v-expand-transition>
    </v-card>
</template>

<script lang="ts">
import Vue from 'vue';
import { DataTableHeader } from 'vuetify';
import SeverityChip from '@/components/SeverityChip.vue';
import { Policy } from '../models';
import YamlDialog from './YamlDialog.vue';

type Item = Policy & { type: string }

type Data = { open: boolean; search: string; expanded: string[] }
type Computed = { headers: DataTableHeader[]; showNamespace: boolean; items: Item[] }
type Props = { title: string; policies: Policy[] }
type Methods = { types(policy: Policy): string }

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { SeverityChip, YamlDialog },
  props: {
    title: { type: String, required: true },
    policies: { type: Array, required: true },
  },
  data: () => ({ open: true, search: '', expanded: [] }),
  computed: {
    showNamespace(): boolean {
      return this.policies.some((p) => !!p.namespace);
    },
    items() {
      return this.policies.map<Item>((p) => ({ ...p, type: this.types(p) }));
    },
    headers(): DataTableHeader[] {
      return [
        { text: 'Kind', value: 'kind' },
        ...(this.showNamespace ? [{ text: 'Namespace', value: 'namespace' }] : []),
        { text: 'Name', value: 'name' },
        { text: 'Rules', value: '' },
        { text: 'Type', value: 'type' },
        { text: 'Background Scan', value: 'background' },
        { text: 'Validation Failure Action', value: 'validationFailureAction' },
        { text: 'Severity', value: 'severity' },
        { text: 'Actions', value: '' },
      ];
    },
  },
  methods: {
    types(policy: Policy): string {
      return policy.rules
        .map((r) => r.type)
        .filter((type, index, self) => self.indexOf(type) === index).join(', ');
    },
  },
});
</script>

<style scoped>
>>> code {
  padding: 16px!important;
}
</style>
