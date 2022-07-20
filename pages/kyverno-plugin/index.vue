<template>
  <loader :loading="loading" :error="error">
    <v-container fluid class="py-6 px-6">
      <v-row>
        <v-col cols="12" md="6">
          <category-chart :policy-groups="policyGroups" />
        </v-col>
        <v-col cols="12" md="6">
          <rule-type-chart :policies="policies" />
        </v-col>
      </v-row>
      <template v-for="(group, category) in policyGroups">
        <v-row v-if="group.length > 0" :key="category">
          <v-col cols="12">
            <policy-table :policies="group" :title="category" />
          </v-col>
        </v-row>
      </template>
    </v-container>
  </loader>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import CategoryChart from '~/policy-reporter-plugins/kyverno/components/CategoryChart.vue'
import PolicyTable from '~/policy-reporter-plugins/kyverno/components/PolicyTable.vue'
import RuleTypeChart from '~/policy-reporter-plugins/kyverno/components/RuleTypeChart.vue'
import { Policy, PolicyGroups } from '~/policy-reporter-plugins/kyverno/types'

type Data = {
  error: Error | null;
  loading: boolean;
  policies: Policy[];
  policyGroups: PolicyGroups;
  interval: any;
}
type Methods = {}
type Props = {}

export default Vue.extend<Data, Methods, {}, Props>({
  name: 'Dashboard',
  components: { CategoryChart, RuleTypeChart, PolicyTable },
  data: () => ({
    error: null,
    policies: [],
    policyGroups: {},
    loading: true,
    interval: null
  }),
  fetch () {
    return this.$kyvernoAPI.policies().then(({ groups, policies }) => {
      this.error = null
      this.policies = policies
      this.policyGroups = groups
    }).catch((error) => {
      this.error = error
      this.policies = []
      this.policyGroups = {}
    }).finally(() => {
      this.loading = false
    })
  },
  computed: mapGetters(['refreshInterval', 'currentCluster']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(() => this.$fetch, refreshInterval)
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
