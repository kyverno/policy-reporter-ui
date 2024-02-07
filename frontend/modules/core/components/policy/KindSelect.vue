<template>
  <v-autocomplete
      multiple
      clearable
      :items="store.kinds.namespaced"
      variant="outlined"
      hide-details
      label="Kinds"
      closable-chips
      :model-value="modelValue"
      @update:model-value="emit('update:modelValue', $event)"
      v-bind="$attrs"
  >
    <template v-slot:selection="{ item, index }">
      <v-chip v-if="index < 2">
        <span>{{ item.title }}</span>
      </v-chip>
      <span v-if="index === 2" class="text-caption align-self-center">
        (+{{ modelValue.length - 2 }} others)
      </span>
    </template>
  </v-autocomplete>
</template>

<script lang="ts" setup>
const props = defineProps<{ source: string; modelValue: string[] }>();

const { store } = useSourceStore(props.source)

const emit = defineEmits<{ 'update:modelValue': [kinds: string[]] }>()
</script>
