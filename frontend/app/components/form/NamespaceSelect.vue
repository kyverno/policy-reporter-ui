<template>
  <v-autocomplete
      multiple
      clearable
      :items="namespaces"
      variant="outlined"
      hide-details
      label="Namespaces"
      closable-chips
      v-model="selected"
      v-bind="$attrs"
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

const namespaces = useNamespacesInjection()

const emit = defineEmits<{ 'update:modelValue': [namespaces: string[]] }>()

watch(selected, (current) => {
  emit('update:modelValue', current)
});
</script>
