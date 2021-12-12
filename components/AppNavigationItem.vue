<template>
  <v-list-item :to="route">
    <slot name="icon">
      <v-list-item-icon v-if="icon">
        <v-icon>{{ icon }}</v-icon>
      </v-list-item-icon>
    </slot>

    <v-list-item-content>
      <v-list-item-title>{{ title }}</v-list-item-title>
    </v-list-item-content>
  </v-list-item>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  props: {
    icon: { type: String, default: null },
    route: { type: String, required: true },
    title: { type: String, required: true }
  }
})
</script>
