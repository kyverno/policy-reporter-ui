<template>
  <v-container v-if="policy" fluid class="py-6 px-6">
    <v-row>
      <v-col>
        <v-card width="100%">
          <v-card-title>
            {{ policy.kind }} "{{ policy.name }}"
            <v-spacer />
            <v-chip v-if="policy.background" class="mb-1" color="green" text-color="white">
              Background scan enabled
            </v-chip>
          </v-card-title>
          <template v-if="policy.description">
            <v-divider />
            <v-card-text>
              {{ policy.description }}
            </v-card-text>
          </template>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="3">
        <policy-details :policy="policy" />
      </v-col>
      <v-col cols="3">
        <v-card style="height: 100%;">
          <v-toolbar flat>
            <v-toolbar-title>Configured Autogen Controller</v-toolbar-title>
          </v-toolbar>
          <v-list>
            <template v-for="controller in policy.autogenControllers">
              <v-divider :key="controller + '_divider'" />
              <v-list-item :key="controller">
                <v-list-item-content>
                  <v-list-item-title>
                    {{ controller }}
                  </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </template>
          </v-list>
        </v-card>
      </v-col>
      <v-col cols="3">
        <v-card min-height="300" style="height: 100%;">
          <v-card-title class="pb-0">
            Passed Results
          </v-card-title>
          <v-card-text
            style="height: calc(100% - 48px); font-size: 9rem !important;"
            :class="['text-center text-h1 d-flex justify-center align-center', `success--text`]"
          >
            {{ passCount }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="3">
        <v-card min-height="300" style="height: 100%;">
          <v-card-title class="pb-0">
            Failed Results
          </v-card-title>
          <v-card-text
            style="height: calc(100% - 48px); font-size: 9rem !important;"
            :class="['text-center text-h1 d-flex justify-center align-center', `error--text`]"
          >
            {{ failCount }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <rules-table :policy="policy" />
      </v-col>
    </v-row>

    <policy-report-table v-if="validations" status="fail" :filter="filter" />
    <policy-report-table v-if="validations" status="pass" :filter="filter" />

    <cluster-policy-report-table v-if="validations" status="fail" :filter="filter" />
    <cluster-policy-report-table v-if="validations" status="pass" :filter="filter" />
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-toolbar flat>
            <v-toolbar-title>
              YAML File
            </v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <yaml-content>
              {{ policy.content }}
            </yaml-content>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-btn
      height="40"
      fixed
      style="z-index: 999; right: 300px; top: 12px;"
      depressed
      dark
      color="grey darken-1"
      @click="$router.push('/kyverno-plugin')"
    >
      Back
    </v-btn>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import { Cluster, Filter, Status } from '~/policy-reporter-plugins/core/types'
import PolicyDetails from '~/policy-reporter-plugins/kyverno/components/PolicyDetails.vue'
import RulesTable from '~/policy-reporter-plugins/kyverno/components/RulesTable.vue'
import YamlContent from '~/policy-reporter-plugins/kyverno/components/YamlContent.vue'
import { Policy, RuleType } from '~/policy-reporter-plugins/kyverno/types'

type Data = {
  policy: Policy | null;
  failCount: number;
  passCount: number;
  interval: any;
}
type Methods = {}
type Props = {}
type Computed = {
  validations: boolean;
  filter: Filter;
  currentCluster?: Cluster;
}

export default Vue.extend<Data, Methods, Computed, Props>({
  name: 'Details',
  components: { RulesTable, PolicyDetails, YamlContent },
  data: () => ({
    interval: null,
    failCount: 0,
    passCount: 0,
    policy: null
  }),
  async fetch () {
    const policy = (await this.$kyvernoAPI.policies().then(({ policies }) => policies)).find((p: Policy) => p.uid === this.$route.params.uid) as Policy

    const namespacedCount = await this.$coreAPI.namespacedStatusCount({ status: [Status.FAIL, Status.PASS], policies: [policy.name] })

    const counter = namespacedCount.reduce<{ [Status.FAIL]: number, [Status.PASS]: number }>((acc, count) => {
      acc[count.status as Status.PASS | Status.FAIL] = count.items.reduce<number>((sum, item) => { return sum + item.count }, 0)

      return acc
    }, { fail: 0, pass: 0 })

    const statusCount = await this.$coreAPI.statusCount({ status: [Status.FAIL, Status.PASS], policies: [policy.name] })

    statusCount.forEach((count) => {
      counter[count.status as Status.PASS | Status.FAIL] += count.count
    })

    this.failCount = counter.fail
    this.passCount = counter.pass
    this.policy = policy
  },
  computed: {
    filter () {
      return { policies: [(this.policy as Policy).name] }
    },
    validations () {
      if (!this.policy) { return false }

      return this.policy.rules.some(r => r.type === RuleType.VALIDATION)
    },
    ...mapGetters(['refreshInterval', 'currentCluster'])
  },
  watch: {
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
    currentCluster: '$fetch'
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>

<style scoped>
>>> code {
  padding: 16px !important;
}
</style>
