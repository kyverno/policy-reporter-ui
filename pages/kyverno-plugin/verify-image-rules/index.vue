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
  asyncData ({ $kyvernoAPI }) {
    return $kyvernoAPI.verifyImageRules().then(rules => ({
      rules,
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
          this.$kyvernoAPI.verifyImageRules().then((rules) => {
            this.rules = rules
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
