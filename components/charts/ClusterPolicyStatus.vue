<template>
  <v-card height="300px">
    <v-card-title class="pb-0">
      {{ statusText }} ClusterPolicies
    </v-card-title>
    <v-card-text :style="`height: calc(100% - 48px); font-size: ${size}rem !important; color: ${color}`" :class="['text-center text-h1 d-flex justify-center align-center']">
      <span v-if="!waiting">{{ count }}</span>
      <v-progress-circular v-else indeterminate size="100" />
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapDarkStatus, mapStatus, mapStatusText } from '~/policy-reporter-plugins/core/mapper'
import { Status } from '~/policy-reporter-plugins/core/types'

export default Vue.extend<{ waiting: boolean }, {}, { statusText: string; color: string; size: number; }, { count: number; status: Status; }>({
  name: 'ClusterPolicyStatus',
  props: {
    status: { required: true, type: String as Vue.PropType<Status> },
    count: { required: true, type: Number }
  },
  data: () => ({
    waiting: true
  }),
  computed: {
    statusText () {
      return mapStatusText(this.status)
    },
    color () {
      if (this.$vuetify.theme.dark) {
        return mapDarkStatus(this.status)
      }

      return mapStatus(this.status)
    },
    size () {
      const counterLength = this.count.toString().length

      if (counterLength <= 3) { return 10 }

      return 10 - ((counterLength - 2))
    }
  },
  created () {
    setTimeout(() => { this.waiting = false }, 500)
  }
})
</script>
