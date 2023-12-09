<template>
  <v-autocomplete
    multiple
    clearable
    density="compact"
    :items="items"
    variant="outlined"
    hide-details
    :label="label"
    closable-chips
    :model-value="selected"
    @update:model-value="input"
    v-bind="$attrs"
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
const props = defineProps({
  label: { type: String, required: true },
  queryName: { type: String, required: true },
  items: { type: Array, required: true },
});

const emit = defineEmits(["itemsLoaded"]);

const selected = ref<string[]>([]);

const input = defineRouteQuery(props.queryName, selected);

selected.value = selected.value.filter((s) => props.items.includes(s));

watch(selected, (current) => {
  emit("itemsLoaded", current);
});
</script>
