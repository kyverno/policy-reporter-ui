<template>
  <loader :loading="loading" :error="error">
    <v-container fluid class="pt-6 px-6">
      <v-row>
        <v-col>
          <v-card>
            <v-card-title class="text-h4" style="text-transform: capitalize;">
              {{ source }} Results
            </v-card-title>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card elevation="1">
            <v-container fluid>
              <v-row>
                <v-col cols="6" class="d-inline-block">
                  <policy-autocomplete :source="source" @input="groupings.policies = $event" />
                </v-col>
                <v-col cols="6">
                  <kind-autocomplete :source="source" />
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <category-autocomplete :source="source" @input="groupings.categories = $event" />
                </v-col>
                <v-col cols="6">
                  <severity-autocomplete />
                </v-col>
              </v-row>
            </v-container>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col v-if="counters['fail']" :cols="sizes.col" :sm="sizes.sm" :md="sizes.md">
          <cluster-policy-status :count="counters['fail']" status="fail" />
        </v-col>
        <v-col v-if="counters['pass']" :cols="sizes.col" :sm="sizes.sm" :md="sizes.md">
          <cluster-policy-status :count="counters['pass']" status="pass" />
        </v-col>
        <v-col v-if="counters['warn']" :cols="sizes.col" :sm="sizes.sm" :md="sizes.md">
          <cluster-policy-status :count="counters['warn']" status="warn" />
        </v-col>
        <v-col v-if="counters['error']" :cols="sizes.col" :sm="sizes.sm" :md="sizes.md">
          <cluster-policy-status :count="counters['error']" status="error" />
        </v-col>
        <v-col v-if="counters['skip']" :cols="sizes.col" :sm="sizes.sm" :md="sizes.md">
          <cluster-policy-status :count="counters['skip']" status="skip" />
        </v-col>
      </v-row>

      <v-row>
        <v-col>
          <v-card elevation="1">
            <v-container fluid>
              <v-row>
                <v-col cols="4" class="d-inline-block">
                  <view-select v-model="groupBy" />
                </v-col>
              </v-row>
            </v-container>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    <v-container v-show="show" fluid class="px-6">
      <template v-if="groupBy === 'status'">
        <cluster-policy-report-table status="error" :filter="{ sources: [source] }" />
        <cluster-policy-report-table status="fail" :filter="{ sources: [source] }" />
        <cluster-policy-report-table status="warn" :filter="{ sources: [source] }" />
        <cluster-policy-report-table status="pass" :filter="{ sources: [source] }" />
        <cluster-policy-report-table status="skip" :filter="{ sources: [source] }" />
      </template>

      <template v-for="value in groupings[groupBy]" v-else>
        <cluster-policy-report-table :key="value" :filter="{ [groupBy]: [value], sources: [source] }" :title-text="value" />
      </template>
    </v-container>
    <v-container v-if="!show">
      <v-row>
        <v-col style="height: 300px;" class="text-center d-flex justify-center align-center">
          <v-progress-circular indeterminate size="150" width="10" color="grey darken-2" />
        </v-col>
      </v-row>
    </v-container>
  </loader>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { boxSizes } from '~/helper/layouthHelper'
import { createCounters, createStatusList, Status } from '~/policy-reporter-plugins/core/types'
import { Policy } from '~/policy-reporter-plugins/kyverno/types'

type Data = {
  error: Error | null;
  show: boolean;
  loading: boolean;
  interval: any;
  counters: { [status in Status]: number };
  groupBy: 'status' | 'policies' | 'categories' | 'rules'
  groupings: {
    policies: Policy[];
    categories: string[];
    status: Status[];
    rules: string[];
  };
}

type Methods = {}
type Computed = { source: string, sizes: { sm: number, md: number, col: number} }
type Props = {}

export default Vue.extend<Data, Methods, Computed, Props>({
  data: () => ({
    error: null,
    show: true,
    loading: true,
    interval: null,
    groupBy: 'status',
    groupings: {
      policies: [],
      rules: [],
      categories: [],
      status: createStatusList()
    },
    counters: createCounters()
  }),
  async fetch () {
    try {
      const [statusCount, rules] = await Promise.all([
        this.$coreAPI.statusCount({ ...this.$route.query, sources: [this.source] }),
        this.$coreAPI.clusterRules(this.source)
      ])

      this.error = null

      statusCount.forEach((item) => {
        this.counters[item.status] = item.count
      })

      this.groupings.rules = rules
    } catch (error) {
      this.groupings.rules = []
      this.counters = createCounters()

      this.error = error as Error
    }

    this.loading = false
  },
  computed: {
    ...mapGetters(['refreshInterval', 'currentCluster']),
    source () {
      return this.$route.params.source
    },
    sizes (): { sm: number, md: number, col: number } {
      return boxSizes(this.counters)
    }
  },
  watch: {
    '$route.query': {
      deep: true,
      immediate: true,
      handler: '$fetch'
    },
    groupBy () {
      this.show = false
      setTimeout(() => {
        this.show = true
      }, 700)
    },
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        if (!refreshInterval) {
          this.interval = null
          return
        }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    },
    currentCluster () {
      this.loading = true
      this.$fetch()
    }
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
