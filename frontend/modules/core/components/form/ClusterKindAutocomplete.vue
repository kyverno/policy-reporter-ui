<template>
  <v-autocomplete
    multiple
    clearable
    density="compact"
    :items="store.kinds.cluster"
    variant="outlined"
    hide-details
    label="Cluster Kinds"
    closable-chips
    :model-value="selected"
    @update:model-value="input"
    v-bind="$attrs"
    v-if="store.kinds.cluster.length"
  >
    <template v-slot:selection="{ item, index }">
      <v-chip v-if="index < 2">
        <span>{{ item.title }}</span>
      </v-chip>
      <span v-if="index === 2" class="text-caption align-self-center">
        (+{{ selected.length - 2 }} others)
      </span>
    </template>
  </v-autocomplete>
</template>

<script lang="ts" setup>
const props = defineProps<{ source?: string; modelValue: string[] }>();

const selected = ref<string[]>(props.modelValue);

const { store } = useSourceStore(props.source)

const input = defineRouteQuery('cluster-kinds', selected);

watch(store.kinds.cluster, (current) => {
  input(selected.value.filter((s) => current.includes(s)));
});

const emit = defineEmits<{ 'update:modelValue': [kinds: string[]] }>()

watch(selected, (current) => {
  emit('update:modelValue', current)
}, { immediate: true });
</script>
