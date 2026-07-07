<template>
  <v-text-field
      :model-value="modelValue"
      type="search"
      @update:model-value="debounced(() => emit('update:modelValue', $event))"
      label="Search"
      hide-details
      variant="outlined"
      density="compact"
      v-bind="$attrs"
      clearable
  />
</template>

<script lang="ts" setup>
defineProps({
  modelValue: { type: String, required: false, default: undefined },
})

const debounced = useDebounce(600)

const emit = defineEmits(['update:modelValue'])
</script>
