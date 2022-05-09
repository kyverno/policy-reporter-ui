<template>
  <div>
    <v-container fluid class="pt-6 px-6">
      <v-row>
        <v-col>
          <v-card elevation="1">
            <v-container fluid>
              <v-row>
                <v-col cols="6" class="d-inline-block">
                  <policy-autocomplete @input="groupings.policies = $event" />
                </v-col>
                <v-col cols="6">
                  <kind-autocomplete />
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <category-autocomplete @input="groupings.categories = $event" />
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
        <v-col cols="12" sm="6" md="3">
          <cluster-policy-status :count="counters['fail']" status="fail" />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <cluster-policy-status :count="counters['pass']" status="pass" />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <cluster-policy-status :count="counters['warn']" status="warn" />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <cluster-policy-status :count="counters['error']" status="error" />
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
        <cluster-policy-report-table status="error" />
        <cluster-policy-report-table status="fail" />
        <cluster-policy-report-table status="warn" />
        <cluster-policy-report-table status="pass" />
        <cluster-policy-report-table status="skip" />
      </template>

      <template v-for="value in groupings[groupBy]" v-else>
        <cluster-policy-report-table :key="value" :filter="{ [groupBy]: [value] }" :title-text="value" />
      </template>
    </v-container>
    <v-container v-if="!show">
      <v-row>
        <v-col style="height: 300px;" class="text-center d-flex justify-center align-center">
          <v-progress-circular indeterminate size="150" width="10" color="grey darken-2" />
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { Status } from '~/policy-reporter-plugins/core/types'
import { Policy } from '~/policy-reporter-plugins/kyverno/types'

type Data = {
  show: boolean;
  interval: any;
  counters: { [status in Status]: number };
  groupBy: 'status' | 'policies' | 'categories' | 'rules'
  groupings: {
    policies: Policy[];
    rules: string[];
    categories: string[];
    status: Status[];
  };
}

type Methods = {}
type Computed = {}
type Props = {}

export default Vue.extend<Data, Methods, Computed, Props>({
  data: () => ({
    show: true,
    interval: null,
    groupBy: 'status',
    groupings: {
      policies: [],
      categories: [],
      rules: [],
      status: [
        Status.FAIL,
        Status.PASS,
        Status.WARN,
        Status.ERROR,
        Status.SKIP
      ]
    },
    counters: {
      [Status.SKIP]: 0,
      [Status.PASS]: 0,
      [Status.WARN]: 0,
      [Status.FAIL]: 0,
      [Status.ERROR]: 0
    }
  }),
  async fetch () {
    const [statusCount, rules] = await Promise.all([
      this.$coreAPI.statusCount(this.$route.query),
      this.$coreAPI.clusterRules()
    ])

    statusCount.forEach((item) => {
      this.counters[item.status] = item.count
    })

    this.groupings.rules = rules
  },
  computed: mapGetters(['refreshInterval']),
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

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
