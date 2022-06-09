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
                <search-field v-model="search" />
              </div>
            </v-card-title>
            <v-divider />
            <v-data-table
              :items="results"
              :server-items-length="count"
              :headers="tableHeaders"
              :options.sync="options"
              :expanded.sync="expanded"
              item-key="id"
            >
              <template #item="{ item, expand, isExpanded }">
                <tr style="cursor: pointer;" @click="expand(!isExpanded)">
                  <td v-if="showResources">
                    <span>{{ item.kind }}</span>
                  </td>
                  <td v-if="showResources">
                    <span>{{ item.name }}</span>
                  </td>
                  <td>
                    <v-chip class="grey-background" label @click="search = item.policy">
                      {{ item.policy }}
                    </v-chip>
                  </td>
                  <td>
                    <v-chip v-if="item.rule" class="grey-background" label @click="search = item.rule">
                      {{ item.rule }}
                    </v-chip>
                  </td>
                  <td>
                    <severity-chip v-if="item.severity" :severity="item.severity" label @click="search = item.severity" />
                  </td>
                  <td>
                    <status-chip :status="item.status" @click="search = item.status" />
                  </td>
                </tr>
              </template>
              <template #expanded-item="{ headers, item }">
                <tr class="table-expand-text">
                  <td :colspan="headers.length" class="py-3">
                    <div v-if="item.hasProps">
                      <v-card v-if="item.message" flat>
                        <v-alert dense type="info" outlined class="rounded" flat>
                          {{ item.message }}
                        </v-alert>
                      </v-card>
                      <div :class="{ 'mt-4': item.message }">
                        <template v-for="(value, label) in item.chips">
                          <property-chip :key="label" :label="label" :value="value" />
                        </template>
                        <template v-for="(value, label) in item.cards">
                          <property-card :key="label" :label="label" :value="value" />
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
import { DataOptions, DataTableHeader } from 'vuetify'
import { mapStatusText } from '~/policy-reporter-plugins/core/mapper'
import { Dictionary, Filter, ListResult, Pagination, Status } from '~/policy-reporter-plugins/core/types'
import { sortByKeys } from '~/helper/layouthHelper'

type Result = ListResult & { chips: Dictionary, cards: Dictionary, hasProps: boolean }

type Data = { open: boolean; search: string; expanded: string[], results: Result[], interval: any; options: DataOptions, count: number }
type Computed = { tableHeaders: DataTableHeader[]; showResources: boolean, title: string }
type Props = { status: Status | null; filter: Filter; titleText: string }

export default Vue.extend<Data, {}, Computed, Props>({
  props: {
    status: { type: String as Vue.PropType<Status>, default: null },
    filter: { type: Object, default: () => ({}) },
    titleText: { type: String, default: '' }
  },
  data: () => ({
    open: true,
    search: '',
    expanded: [],
    results: [],
    count: 0,
    interval: null,
    options: {
      itemsPerPage: 10,
      page: 1,
      sortDesc: [],
      sortBy: [],
      groupBy: [],
      groupDesc: [],
      multiSort: false,
      mustSort: false
    }
  }),
  fetch () {
    const filter = { ...(this.status ? { status: [this.status] } : {}), ...this.$route.query, ...this.filter }
    if (this.search) {
      filter.search = this.search
    }

    const pagination: Pagination = {
      page: this.options.page,
      offset: this.options.itemsPerPage
    }
    return this.$coreAPI.results(filter, pagination).then(({ items, count }) => {
      this.results = items.map(({ properties, ...result }) => {
        const chips: Dictionary = {}
        const cards: Dictionary = {}
        let hasProps: boolean = false

        for (const prop in properties) {
          if (properties[prop].length > 75) {
            cards[prop] = properties[prop]
          } else {
            chips[prop] = properties[prop]
          }
          hasProps = true
        }

        return {
          ...result,
          properties: {},
          cards: sortByKeys(cards),
          chips: sortByKeys(chips),
          hasProps
        }
      })
      this.count = count
    })
  },
  computed: {
    title () {
      if (this.status) {
        return `${mapStatusText(this.status)} ClusterPolicy Results`
      }

      return this.titleText
    },
    showResources (): boolean {
      return this.results.some(item => !!item.kind)
    },
    tableHeaders (): DataTableHeader[] {
      const resourceFileds = this.showResources
        ? [
            { text: 'Kind', value: 'kind' },
            { text: 'Name', value: 'name' }
          ]
        : []

      return [
        ...resourceFileds,
        { text: 'Policy', value: 'policy' },
        { text: 'Rule', value: 'rule' },
        { text: 'Severity', value: 'severity' },
        { text: 'Status', value: 'status' }
      ]
    },
    ...mapGetters(['refreshInterval'])
  },
  watch: {
    '$route.query': {
      deep: true,
      handler: '$fetch'
    },
    search: {
      deep: true,
      handler: '$fetch'
    },
    options: {
      handler: '$fetch',
      deep: true
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
