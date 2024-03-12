<template>
  <v-select
      multiple
      clearable
      density="compact"
      :items="items"
      variant="outlined"
      hide-details
      label="Results"
      closable-chips
      :model-value="modelValue"
      @update:model-value="input"
      v-bind="$attrs"
      style="min-width: 200px;"
  >
    <template v-slot:selection="{ item, index }">
      <v-chip v-if="index < 2">
        <span>{{ item.title }}</span>
      </v-chip>
      <span v-if="index === 2" class="text-caption align-self-center">
        (+{{ modelValue.length - 2 }} others)
      </span>
    </template>
  </v-select>
</template>

<script lang="ts" setup>
import {Status} from "../../types";

const props = defineProps<{ source?: string; modelValue: string[] }>();

const items: Status[] = [
  Status.PASS,
  Status.FAIL,
  Status.WARN,
  Status.SKIP,
]

const emit = defineEmits<{ 'update:modelValue': [status: string[]] }>()

const input = (current) => emit('update:modelValue', current);

</script>
