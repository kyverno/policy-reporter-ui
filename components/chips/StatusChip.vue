<template>
  <v-chip :color="color" :dark="dark" v-bind="$attrs" v-on="$listeners">
    {{ status }}
  </v-chip>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapDarkStatus, mapStatus } from '~/policy-reporter-plugins/core/mapper'
import { Status } from '~/policy-reporter-plugins/core/types'

export default Vue.extend<{}, {}, { color: string; dark: boolean }, { status: Status }>({
  inheritAttrs: false,
  props: {
    status: { type: String as Vue.PropType<Status>, required: true }
  },
  computed: {
    color () {
      return this.$vuetify.theme.dark ? mapDarkStatus(this.status) : mapStatus(this.status)
    },
    dark () {
      if ([Status.PASS, Status.WARN, Status.ERROR, Status.FAIL].includes(this.status)) {
        return true
      }

      return false
    }
  }
})
</script>
