<template>
  <loader :loading="loading" :error="$fetchState.error">
    <v-container fluid class="py-6 px-6">
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
            <v-data-table
              :search="search"
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
                  <td class="py-1 pl-3 pr-1" style="white-space: nowrap; width: 200px;">
                    {{ formatTime(item.creationTimestamp) }}
                  </td>
                  <td class="py-1 px-1" style="white-space: nowrap;">
                    [{{ item.priority.toUpperCase() }}]
                  </td>
                  <td class="py-1 px-1" style="white-space: nowrap;">
                    {{ item.kind }}
                  </td>
                  <td class="py-1 px-1" style="white-space: nowrap;">
                    [{{ item.namespace }}]
                  </td>
                  <td class="py-1 px-1" style="white-space: nowrap;">
                    {{ item.name }}
                  </td>
                  <td class="py-1 px-1" style="white-space: nowrap; overflow: scroll;">
                    {{ item.message }}
                  </td>
                </tr>
              </template>
            </v-data-table>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </loader>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { DataTableHeader } from 'vuetify'

type Data = {
  loading: boolean;
  search: string;
  expanded: any[];
  headers: DataTableHeader[];
  items: any[];
  interval: any;
}

type Methods = {
  formatTime(date: string): string;
}

type Computed = {
  refreshInterval: number;
}

export default Vue.extend<Data, Methods, Computed, {}>({
  name: 'Log',
  components: {
  },
  data: () => ({
    loading: true,
    interval: null,
    items: [],
    search: '',
    expanded: ['creationTimestamp', 1],
    headers: [
      { text: 'dateime', value: 'creationTimestamp', sortable: false },
      { text: 'priority', value: 'priority', sortable: false },
      { text: 'kind', value: 'kind', sortable: false },
      { text: 'namespace', value: 'kind', sortable: false },
      { text: 'name', value: 'name', sortable: false },
      { text: 'message', value: 'message', sortable: false }
    ]
  }),
  fetch () {
    return this.$coreAPI.logs().then((logs) => {
      this.items = logs.map(({ resource, ...log }) => ({ ...resource, ...log }))
    }).finally(() => {
      this.loading = false
    })
  },
  computed: mapGetters(['refreshInterval']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        if (!refreshInterval) {
          this.interval = null
          this.$fetch()
          return
        }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  },
  methods: {
    formatTime (date: string): string {
      return new Date(date).toISOString().replace('T', ' ')
    }
  }
})
</script>
