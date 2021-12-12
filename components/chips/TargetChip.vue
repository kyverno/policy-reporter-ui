<template>
  <v-chip
    :key="target.name"
    class="ml-3"
    outlined
    label
    style="height: 40px;"
    :title="buildTitle(target)"
  >
    <v-avatar left>
      <v-icon :color="target.minimumPriority | mapPriority">
        mdi-target-variant
      </v-icon>
    </v-avatar>
    {{ target.name }}
  </v-chip>
</template>

<script lang="ts">
import Vue from 'vue'
import { Target } from '~/policy-reporter-plugins/core/types'

export default Vue.extend({
  props: {
    target: { type: Object as Vue.PropType<Target>, required: true }
  },
  methods: {
    buildTitle (target: Target) {
      let title = `minimum priority: ${target.minimumPriority}`

      if (target.sources && target.sources.length > 0) {
        title += `\nsource filter: ${target.sources.join(', ')}`
      }

      return title
    }
  }
})
</script>
