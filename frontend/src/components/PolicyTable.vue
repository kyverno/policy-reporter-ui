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
              :sort-by="showResources ? ['resource.namespace', 'resource.name', 'policy', 'rule'] : ['policy', 'rule']"
              :expanded.sync="expanded"
              item-key="id"
              >
            <template #item="{ item, expand, isExpanded }">
                <tr @click="expand(!isExpanded)" style="cursor: pointer">
                  <td v-if="showResources"><span v-if="item.resource">{{ item.resource.namespace }}</span></td>
                  <td v-if="showResources"><span v-if="item.resource">{{ item.resource.kind }}</span></td>
                  <td v-if="showResources"><span v-if="item.resource">{{ item.resource.name }}</span></td>
                  <td>
                    <v-chip class="grey lighten-2" label @click.stop="search = item.policy">
                      {{ item.policy }}
                    </v-chip>
                  </td>
                  <td>
                    <v-chip class="grey lighten-2" label @click.stop="search = item.rule">
                      {{ item.rule }}
                    </v-chip>
                  </td>
                  <td>
                    <severity-chip :severity="item.severity" label @click.stop="search = item.severity" v-if="item.severity" />
                  </td>
                  <td>
                    <status-chip @click.stop="search = item.status" :status="item.status" />
                  </td>
                </tr>
            </template>
            <template #expanded-item="{ headers, item }">
              <tr class="grey lighten-4">
                <td :colspan="headers.length" class="py-3">
                  <div v-if="item.properties">
                    <v-card flat>
                      <v-alert type="info" outlined class="rounded" flat>
                        {{ item.message }}
                      </v-alert>
                    </v-card>
                    <div class="mt-4">
                      <template v-for="(value, label) in item.properties">
                        <property-chip :key="label" :label="label" :value="value" v-if="value.length <= 100" />
                        <property-card :key="label" :label="label" :value="value" v-else />
                      </template>
                    </div>
                  </div>
                  <div v-else>
                  {{ item.message }}
                  </div>
                </td>
              </tr>
            </template>
            </v-data-table>
        </div>

        </v-expand-transition>
    </v-card>
</template>

<script lang="ts">
import { Result } from '@/models';
import Vue from 'vue';
import { DataTableHeader } from 'vuetify';
import PropertyCard from './PropertyCard.vue';
import PropertyChip from './PropertyChip.vue';
import SeverityChip from './SeverityChip.vue';
import StatusChip from './StatusChip.vue';

type Data = { open: boolean; search: string; expanded: string[] }
type Computed = { headers: DataTableHeader[]; items: Item[]; showResources: boolean }
type Props = { title: string; results: Result[] }
type Methods = {}

type Item = Result & { id: string }

export default Vue.extend<Data, Methods, Computed, Props>({
  components: {
    StatusChip, SeverityChip, PropertyChip, PropertyCard,
  },
  props: {
    title: { type: String, required: true },
    results: { type: Array, required: true },
  },
  data: () => ({ open: true, search: '', expanded: [] }),
  computed: {
    items(): Item[] {
      return this.results.map((result: Result) => {
        const properties = { ...(result.properties || {}) };

        const sorted = Object.entries(properties);

        sorted.sort(([ak, av], [bk, bv]) => (ak.length + av.length) - (bk.length + bv.length));

        return { ...result, id: result.policy + result.rule + result?.resource?.uid, properties: Object.fromEntries(sorted) };
      });
    },
    showResources(): boolean {
      return this.results.some((item) => !!item.resource);
    },
    headers(): DataTableHeader[] {
      const resourceFileds = this.showResources ? [
        { text: 'Namespace', value: 'resource.namespace' },
        { text: 'Kind', value: 'resource.kind' },
        { text: 'Name', value: 'resource.name' },
      ] : [];

      return [
        ...resourceFileds,
        { text: 'Policy', value: 'policy' },
        { text: 'Rule', value: 'rule' },
        { text: 'Severity', value: 'severity' },
        { text: 'Status', value: 'status' },
      ];
    },
  },
});
</script>
