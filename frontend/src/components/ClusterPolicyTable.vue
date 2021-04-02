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
                />
              </div>
            </v-card-title>
            <v-divider />
            <v-data-table :items="results" :headers="headers" :items-per-page="10" :search="search" :sort-by="['resource.name', 'resource.kind']">
            <template #item="{ item }">
                <tr>
                  <td>{{ item.resource.kind }}</td>
                  <td>{{ item.resource.name }}</td>
                  <td>
                    <v-chip class="grey lighten-2" label @click="search = item.policy">
                      {{ item.policy }}
                    </v-chip>
                  </td>
                  <td>
                    <v-chip class="grey lighten-2" label @click="search = item.rule">
                      {{ item.rule }}
                    </v-chip>
                  </td>
                  <td>
                    <severity-chip :severity="item.severity" label @click="search = item.severity" v-if="item.severity" />
                  </td>
                  <td>
                    <status-chip @click="search = item.status" :status="item.status" />
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
import SeverityChip from './SeverityChip.vue';
import StatusChip from './StatusChip.vue';

type Data = { open: boolean; search: string }
type Computed = { headers: DataTableHeader[]; severityConfigured: boolean }
type Props = { title: string; results: Result[] }

export default Vue.extend<Data, {}, Computed, Props>({
  components: { StatusChip, SeverityChip },
  props: {
    title: { type: String, required: true },
    results: { type: Array, required: true },
  },
  data: () => ({ open: true, search: '' }),
  computed: {
    headers(): DataTableHeader[] {
      return [
        { text: 'Kind', value: 'resource.kind' },
        { text: 'Name', value: 'resource.name' },
        { text: 'Policy', value: 'policy' },
        { text: 'Rule', value: 'rule' },
        { text: 'Severity', value: 'severity' },
        { text: 'Status', value: 'status' },
      ];
    },
    severityConfigured() {
      return this.results.some((result) => !!result.severity);
    },
  },
});
</script>
