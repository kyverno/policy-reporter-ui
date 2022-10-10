<template>
  <v-chip
    dark
    label
    small
    :class="`my-1 mr-2 pr-0 py-0`"
    color="info"
    style="border: 1px solid #000;"
    v-on="$listeners"
  >
    {{ label }}
    <div class="property-text-bg rounded ml-2 px-2 py-2 text-body-2" style="font-size: 0.7rem !important; line-height: 0.9rem !important;">
      <a v-if="isURL(value)" :href="value" target="_blank">{{ value }}</a>
      <span v-else v-html="value" />
    </div>
  </v-chip>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  name: 'PropertyChip',
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true }
  },
  methods: {
    isURL (value: string): boolean {
      let url

      try {
        url = new URL(value)
      } catch (_) {
        return false
      }

      return url.protocol === 'http:' || url.protocol === 'https:'
    }
  }
})
</script>
