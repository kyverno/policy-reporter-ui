<template>
  <v-chip :color="color" dark v-bind="$attrs" v-on="$listeners">
    {{ status }}
  </v-chip>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapDarkStatus, mapStatus } from '~/policy-reporter-plugins/core/mapper'
import { Status } from '~/policy-reporter-plugins/core/types'

export default Vue.extend<{}, {}, { color: string; }, { status: Status }>({
  inheritAttrs: false,
  props: {
    status: { type: String as Vue.PropType<Status>, required: true }
  },
  computed: {
    color () {
      return this.$vuetify.theme.dark ? mapDarkStatus(this.status) : mapStatus(this.status)
    }
  }
})
</script>
