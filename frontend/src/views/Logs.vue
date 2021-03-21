<template>
  <v-container fluid class="py-8 px-6">
    <v-row>
        <v-col>
          <v-card>
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
            <v-data-table :search="search"
                          :items-per-page="-1"
                          sort-by="creationTimestamp"
                          sort-desc
                          :items="items"
                          hide-default-footer
                          hide-default-header
                          item-key="creationTimestamp"
                          :headers="headers"
            >
              <template #item="{ item }">
                  <tr>
                    <td class="py-1 pl-3 pr-1" style="white-space:nowrap;width: 200px;">{{ formatTime(item.creationTimestamp) }}</td>
                    <td class="py-1 px-1" style="white-space:nowrap;">[{{ item.priority.toUpperCase() }}]</td>
                    <td class="py-1 px-1" style="white-space:nowrap;">{{ item.kind }}</td>
                    <td class="py-1 px-1" style="white-space:nowrap;">{{ item.name }}</td>
                    <td class="py-1 px-1" style="white-space:nowrap;overflow: scroll">{{ item.message }}</td>
                  </tr>
              </template>
            </v-data-table>
          </v-card>
        </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue';
import { Result } from '@/models';
import { FETCH_LOG } from '@/store';
import { mapState } from 'vuex';
import { DataTableHeader } from 'vuetify';

type Data = {
  search: string;
  expanded: any[];
  headers: DataTableHeader[];
}

type Methods = {
  formatTime(date: string): string;
}

type Computed = {
  log: Result[];
  items: any[];
}

export default Vue.extend<Data, Methods, Computed, {}>({
  components: {
  },
  name: 'Log',
  data: () => ({
    search: '',
    expanded: ['creationTimestamp', 1],
    headers: [
      { text: 'dateime', value: 'creationTimestamp', sortable: false },
      { text: 'priority', value: 'priority', sortable: false },
      { text: 'kind', value: 'kind', sortable: false },
      { text: 'name', value: 'name', sortable: false },
      { text: 'message', value: 'message', sortable: false },
    ],
  }),
  computed: {
    ...mapState(['log']),
    items() {
      return this.log.map(({ resource, ...log }) => ({ ...resource, ...log }));
    },
  },
  methods: {
    formatTime(date: string): string {
      return new Date(date).toISOString().replace('T', ' ');
    },
  },
  created() {
    this.$store.dispatch(FETCH_LOG);

    setInterval(() => {
      this.$store.dispatch(FETCH_LOG);
    }, 1000);
  },
});
</script>
