<template>
  <v-card variant="outlined" v-bind="$attrs" class="bg-white">
    <v-card-title class="bg-indigo"><div class="text-subtitle-2">{{ label }}</div></v-card-title>
    <v-divider />
    <v-card-text class="text-body-2 pa-2">
      <a v-if="isURL(value)" :href="value" target="_blank">{{ value }}</a>
      <span v-else v-html="value" />
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
const props = defineProps<{ label: string; value: string; }>()

const isURL = (value: string): boolean => {
  let url

  try {
    url = new URL(value)
  } catch (_) {
    return false
  }

  return url.protocol === 'http:' || url.protocol === 'https:'
}
</script>
