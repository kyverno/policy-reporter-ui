<template>
  <v-autocomplete
      multiple
      clearable
      density="compact"
      :items="items"
      variant="outlined"
      hide-details
      label="Namespaces"
      closable-chips
      :model-value="selected"
      @update:model-value="input"
      style="min-width: 300px;"
  >
    <template #selection="{ item, index }">
      <v-chip v-if="index <= 1" size="small" label>
        <span>{{ (item as any).title }}</span>
      </v-chip>
      <span v-if="index === 2" class="grey--text caption ml-2 d-inline-flex align-center">
        (+{{ selected.length - 2 }} others)
      </span>
    </template>
  </v-autocomplete>
</template>

<script lang="ts" setup>
import { APIFilter } from "~/modules/core/provider/dashboard";
import { useDebounce } from "~/composables/router";

const props = defineProps<{ modelValue: string[]; items: string[]; }>();

const emit = defineEmits<{ 'update:modelValue': string[] }>();

const filter = inject(APIFilter, ref({}))

const selected = ref<string[]>([]);

const input = defineRouteQuery('namespaces', selected);

watch(() => props.items, (current) => {
  input(selected.value.filter((s) => current.includes(s)));
}, { immediate: true });

const debounced = useDebounce()

watch(selected, (current) => debounced(() => emit("update:modelValue", current)), { immediate: true });
</script>
