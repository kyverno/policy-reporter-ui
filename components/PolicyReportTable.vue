<template>
  <v-row v-if="results.length">
    <v-col cols="12">
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
              :items="results"
              :headers="tableHeaders"
              :items-per-page="10"
              :search="search"
              :sort-by="showResources ? ['namespace', 'name', 'policy', 'rule', 'message'] : ['policy', 'rule', 'message']"
              :expanded.sync="expanded"
              item-key="id"
            >
              <template #item="{ item, expand, isExpanded }">
                <tr style="cursor: pointer;" @click.stop="expand(!isExpanded)">
                  <td v-if="showResources">
                    <span>{{ item.namespace }}</span>
                  </td>
                  <td v-if="showResources">
                    <span>{{ item.kind }}</span>
                  </td>
                  <td v-if="showResources">
                    <span>{{ item.name }}</span>
                  </td>
                  <td>
                    <v-chip class="grey-background" label @click.stop="search = item.policy">
                      {{ item.policy }}
                    </v-chip>
                  </td>
                  <td>
                    <v-chip v-if="item.rule" class="grey-background" label @click.stop="search = item.rule">
                      {{ item.rule }}
                    </v-chip>
                  </td>
                  <td>
                    <severity-chip v-if="item.severity" :severity="item.severity" label @click.stop="search = item.severity" />
                  </td>
                  <td>
                    <status-chip :status="item.status" @click.stop="search = item.status" />
                  </td>
                </tr>
              </template>
              <template #expanded-item="{ headers, item }">
                <tr class="table-expand-text">
                  <td :colspan="headers.length" class="py-3">
                    <div v-if="item.properties && item.properties.length">
                      <v-card flat>
                        <v-alert dense type="info" outlined class="rounded my-0" flat>
                          {{ item.message }}
                        </v-alert>
                      </v-card>
                      <div class="mt-4">
                        <template v-for="(value, label) in item.properties">
                          <property-chip v-if="value.length <= 100" :key="label" :label="label" :value="value" />
                          <property-card v-else :key="label" :label="label" :value="value" />
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
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { DataTableHeader } from 'vuetify'
import { mapStatusText } from '~/policy-reporter-plugins/core/mapper'
import { Filter, ListResult, Status } from '~/policy-reporter-plugins/core/types'

type Data = { open: boolean; search: string; expanded: string[], results: ListResult[]; interval: any }
type Computed = { tableHeaders: DataTableHeader[]; showResources: boolean ; title: string; refreshInterval: number }
type Methods = {}
type Props = { status: Status | null; filter: Filter; titleText: string }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    status: { type: String as Vue.PropType<Status>, default: null },
    filter: { type: Object, default: () => ({}) },
    titleText: { type: String, default: '' }
  },
  data: () => ({ open: true, search: '', expanded: [], results: [], interval: null }),
  fetch () {
    return this.$coreAPI.namespacedResults({ ...(this.status ? { status: [this.status] } : {}), ...this.$route.query, ...this.filter }).then((results) => {
      this.results = results.map(({ properties, ...result }) => ({
        ...result, properties: Object.fromEntries(Object.entries(properties || {}).sort(([, a], [, b]) => a.length - b.length))
      }))
    })
  },
  computed: {
    title () {
      if (this.status) {
        return `${mapStatusText(this.status)} Policy Results`
      }

      return this.titleText
    },
    showResources (): boolean {
      return this.results.some(item => !!item.kind)
    },
    tableHeaders (): DataTableHeader[] {
      const resourceFileds = this.showResources
        ? [
            { text: 'Namespace', value: 'namespace' },
            { text: 'Kind', value: 'kind' },
            { text: 'Name', value: 'name' }
          ]
        : []

      return [
        ...resourceFileds,
        { text: 'Policy', value: 'policy' },
        { text: 'Rule', value: 'rule' },
        { text: 'Severity', value: 'severity', width: 120 },
        { text: 'Status', value: 'status', width: 100 }
      ]
    },
    ...mapGetters(['refreshInterval'])
  },
  watch: {
    '$route.query': {
      deep: true,
      handler: '$fetch'
    },
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
