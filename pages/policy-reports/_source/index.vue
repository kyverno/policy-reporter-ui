<template>
  <loader :loading="loading" :error="$fetchState.error">
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
                <v-col cols="5" class="d-inline-block">
                  <policy-autocomplete :source="source" namespaced @input="groupings.policies = $event" />
                </v-col>
                <v-col cols="4">
                  <kind-autocomplete :source="source" namespaced />
                </v-col>
                <v-col cols="3" class="d-inline-block">
                  <severity-autocomplete />
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="5">
                  <category-autocomplete :source="source" @input="groupings.categories = $event" />
                </v-col>
                <v-col cols="4">
                  <namespace-autocomplete :source="source" />
                </v-col>
              </v-row>
            </v-container>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <policy-status-per-namespace
          status="fail"
          :min-height="minHeight"
          :values="counters['fail']"
          :full-width="!counters['pass'].namespaces.length"
          optional
          @height-change="updateHeight('fail', $event)"
        />
        <policy-status-per-namespace
          status="pass"
          :min-height="minHeight"
          :values="counters['pass']"
          :full-width="!counters['fail'].namespaces.length"
          optional
          @height-change="updateHeight('pass', $event)"
        />
        <policy-status-per-namespace
          status="error"
          :min-height="minHeight"
          :values="counters['error']"
          :full-width="!counters['warn'].namespaces.length"
          optional
          @height-change="updateHeight('error', $event)"
        />
        <policy-status-per-namespace
          status="warn"
          :min-height="minHeight"
          :values="counters['warn']"
          :full-width="!counters['error'].namespaces.length"
          optional
          @height-change="updateHeight('warn', $event)"
        />
        <policy-status-per-namespace
          status="skip"
          :min-height="minHeight"
          :values="counters['skip']"
          :full-width="true"
          optional
          @height-change="updateHeight('skip', $event)"
        />
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
        <policy-report-table status="fail" :filter="{ sources: [source] }" />
        <policy-report-table status="pass" :filter="{ sources: [source] }" />
        <policy-report-table status="warn" :filter="{ sources: [source] }" />
        <policy-report-table status="error" :filter="{ sources: [source] }" />
        <policy-report-table status="skip" :filter="{ sources: [source] }" />
      </template>

      <template v-for="value in groupings[groupBy]" v-else>
        <policy-table :key="value" :filter="{ [groupBy]: [value], sources: [source] }" :title-text="value" />
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
import { Status } from '~/policy-reporter-plugins/core/types'
import { Policy } from '~/policy-reporter-plugins/kyverno/types'

type Data = {
  loading: boolean;
  show: boolean;
  interval: any;
  heights: { [key in Status]: number };
  counters: { [status in Status]: { namespaces: string[]; counts: number[] } };
  groupBy: 'status' | 'policies' | 'categories'
  groupings: {
    policies: Policy[];
    categories: string[];
    status: Status[];
  };
}

type Methods = {
  updateHeight(status: string, height: number): void;
}

type Computed = {
  minHeight: number;
  refreshInterval: number;
  source: string;
}
type Props = {}

export default Vue.extend<Data, Methods, Computed, Props>({
  data: () => ({
    loading: true,
    show: false,
    interval: null,
    groupBy: 'status',
    groupings: {
      policies: [],
      categories: [],
      status: [
        Status.FAIL,
        Status.PASS,
        Status.WARN,
        Status.ERROR,
        Status.SKIP
      ]
    },
    heights: {
      [Status.SKIP]: 0,
      [Status.PASS]: 0,
      [Status.WARN]: 0,
      [Status.FAIL]: 0,
      [Status.ERROR]: 0
    },
    counters: {
      [Status.SKIP]: { namespaces: [], counts: [] },
      [Status.PASS]: { namespaces: [], counts: [] },
      [Status.WARN]: { namespaces: [], counts: [] },
      [Status.FAIL]: { namespaces: [], counts: [] },
      [Status.ERROR]: { namespaces: [], counts: [] }
    }
  }),
  async fetch () {
    const namespacedStatusCount = await this.$coreAPI.namespacedStatusCount({ ...this.$route.query, sources: [this.source] })

    this.counters = namespacedStatusCount.reduce((counters, statusCount) => {
      counters[statusCount.status] = statusCount.items.reduce<{ namespaces: string[]; counts: number[] }>((acc, statusCount) => {
        acc.namespaces.push(statusCount.namespace)
        acc.counts.push(statusCount.count)

        return acc
      }, { namespaces: [], counts: [] })

      return counters
    }, { ...this.counters })

    this.loading = false
  },
  computed: {
    minHeight () {
      return Object.values(this.heights).reduce<number>((acc, height) => {
        if (acc > height) { return acc }

        return height
      }, 0)
    },
    source () {
      return this.$route.params.source
    },
    ...mapGetters(['refreshInterval'])
  },
  watch: {
    '$route.query': {
      deep: true,
      immediate: true,
      handler: '$fetch'
    },
    groupBy: {
      immediate: true,
      handler () {
        this.show = false
        setTimeout(() => {
          this.show = true
        }, 700)
      }
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
  },
  methods: {
    updateHeight (status: string, height: number) {
      this.heights = { ...this.heights, [status as Status]: height }
    }
  }
})
</script>
