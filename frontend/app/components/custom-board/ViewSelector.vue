<template>
  <v-btn-toggle
    :model-value="modelValue"
    mandatory
    density="comfortable"
    variant="outlined"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-btn v-for="view in views" :key="view" :value="view" size="small">
      {{ labels[view] }}
    </v-btn>
  </v-btn-toggle>
</template>

<script setup lang="ts">
type ResultView = 'resources' | 'results' | 'policies'

withDefaults(defineProps<{
  modelValue: ResultView
  views: ResultView[]
}>(), {
  views: () => [],
})

defineEmits<{
  'update:modelValue': [view: ResultView]
}>()

const labels: Record<ResultView, string> = {
  resources: 'Resources',
  results: 'Results',
  policies: 'Policies',
}
</script>
