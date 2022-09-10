<template>
  <loader :loading="loading">
    <v-container fluid class="py-6 px-6">
      <verify-image-table :rules="rules" title="VerifyImages Rules" />
    </v-container>
  </loader>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters } from 'vuex'
import VerifyImageTable from '~/policy-reporter-plugins/kyverno/components/VerifyImageTable.vue'
import { VerifyImageRule } from '~/policy-reporter-plugins/kyverno/types'

type Data = {
  loading: boolean;
  rules: VerifyImageRule[];
  interval: any;
}
type Methods = {}
type Props = {}

export default Vue.extend<Data, Methods, {}, Props>({
  name: 'VerifyImageRules',
  components: { VerifyImageTable },
  data: () => ({
    rules: [],
    interval: null,
    loading: true
  }),
  fetch () {
    return this.$kyvernoAPI.verifyImageRules().then((rules) => {
      this.rules = rules
      this.loading = false
    })
  },
  computed: mapGetters(['refreshInterval', 'currentCluster']),
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
