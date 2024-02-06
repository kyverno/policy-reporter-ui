<template>
  <v-select
    :model-value="modelValue"
    :items="items"
    variant="outlined"
    hide-details
    @update:model-value="value => $emit('update:modelValue', value)"
  />
</template>

<script lang="ts" setup>
defineProps<{ modelValue: string }>()

const items = [
  { title: 'Policy Report', value: 'policy-report' },
  { title: 'Namespace Report', value: 'namespace-report' },
]

</script>
