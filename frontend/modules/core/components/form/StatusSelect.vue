<template>
  <v-select
      multiple
      clearable
      density="compact"
      :items="items"
      variant="outlined"
      hide-details
      label="Status"
      closable-chips
      :model-value="selected"
      @update:model-value="input"
      v-bind="$attrs"
  >
    <template v-slot:selection="{ item, index }">
      <v-chip v-if="index < 2">
        <span>{{ item.title }}</span>
      </v-chip>
      <span v-if="index === 2" class="text-caption align-self-center">
        (+{{ selected.length - 2 }} others)
      </span>
    </template>
  </v-select>
</template>

<script lang="ts" setup>
import {Status} from "../../types";

const props = defineProps<{ source?: string; modelValue?: string[] }>();

const selected = ref<string[]>(props.modelValue || []);

const items: Status[] = [
  Status.PASS,
  Status.FAIL,
  Status.WARN,
  Status.SKIP,
]

const input = defineRouteQuery('status', selected);

const emit = defineEmits<{ 'update:modelValue': [status: string[]] }>()

watch(selected, (current) => {
  emit('update:modelValue', current)
}, { immediate: true });
</script>
