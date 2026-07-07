<template>
  <v-chip variant="flat" color="indigo-lighten-1" :label="true" class="text-label-medium">
    <span class="inline-block mr-1">{{ capilize(label) }}:</span> <a v-if="isURL(value)" :href="value" target="_blank" class="text-white">{{ value }}</a><span v-else>{{ value }}</span>
  </v-chip>
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
