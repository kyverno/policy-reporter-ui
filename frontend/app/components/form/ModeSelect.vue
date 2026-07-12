<template>
  <v-select
    class="mr-2"
    :model-value="modelValue"
    :items="modes"
    variant="outlined"
    hide-details
    density="compact"
    style="min-width: 160px;"
    prepend-inner-icon="mdi-view-compact"
    @update:model-value="select"
    v-if="modes.length > 1"
  />
</template>

<script lang="ts" setup>
import type { Mode } from '~/types/core';

const props = defineProps<{ modelValue: Mode; }>();

const router = useRouter()
const route = useRoute()

const modes: Mode[] = ['detailed', 'compact'];

const emit = defineEmits<{ 'update:modelValue': [mode: Mode] }>()

const select = (mode: Mode) => {
  emit('update:modelValue', mode)
  router.push({ name: route.name as string, query: { ...route.query, mode }, params: route.params })
}

</script>
