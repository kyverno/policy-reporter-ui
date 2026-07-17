<template>
  <v-autocomplete
    multiple
    clearable
    density="compact"
    :items="kinds"
    variant="outlined"
    hide-details
    label="Cluster Kinds"
    closable-chips
    :model-value="selected"
    @update:model-value="input"
    v-bind="$attrs"
    v-if="kinds.length"
  >
    <template v-slot:selection="{ item, index }">
      <v-chip v-if="index < 2">
        <span>{{ item }}</span>
      </v-chip>
      <span v-if="index === 2" class="text-caption align-self-center">
        (+{{ selected.length - 2 }} others)
      </span>
    </template>
  </v-autocomplete>
</template>

<script lang="ts" setup>
const props = defineProps<{ modelValue: string[] }>();

const selected = ref<string[]>(props.modelValue);
const kinds = useClusterKindsInjection()

const input = defineRouteQuery('cluster-kinds', selected);

watch(kinds, (current) => {
  input(selected.value.filter((s) => current.includes(s)));
});

const emit = defineEmits<{ 'update:modelValue': [kinds: string[]] }>()

watch(selected, (current) => {
  emit('update:modelValue', current)
}, { immediate: true });
</script>
