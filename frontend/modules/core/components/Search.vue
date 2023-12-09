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
  />
</template>

<script lang="ts" setup>
import debounce from 'lodash.debounce'

defineProps({
  modelValue: { type: String, required: false, default: undefined },
})

const debounced = debounce((emit: () => void) => { emit() }, 600)

const emit = defineEmits(['update:modelValue'])
</script>
