<template>
  <loader :loading="loading">
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
  layout: 'KyvernoPlugin',
  asyncData ({ $kyvernoAPI }) {
    return $kyvernoAPI.policies().then(({ groups, policies }) => ({
      policies,
      policyGroups: groups,
      interval: null,
      loading: true
    }))
  },
  computed: mapGetters(['refreshInterval']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(() => {
          this.$kyvernoAPI.policies().then(({ groups, policies }) => {
            this.policies = policies
            this.policyGroups = groups
            this.loading = false
          })
        }, refreshInterval)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  }
})
</script>
