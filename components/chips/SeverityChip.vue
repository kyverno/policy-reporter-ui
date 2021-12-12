<template>
  <v-chip :color="color" dark v-bind="$attrs" v-on="$listeners">
    {{ severity }}
  </v-chip>
</template>

<script lang="ts">
import Vue from 'vue'
import { Severity } from '~/policy-reporter-plugins/core/types'

export default Vue.extend<{}, {}, { color: string }, { severity: Severity }>({
  inheritAttrs: false,
  props: {
    severity: { type: String as Vue.PropType<Severity>, required: true }
  },
  computed: {
    color () {
      switch (this.severity) {
        case Severity.LOW:
          return 'info'
        case Severity.MEDIUM:
          return 'warning'
        case Severity.HIGH:
          return 'error'
        default:
          return 'grey'
      }
    }
  }
})
</script>
