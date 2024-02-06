<template>
  <v-autocomplete
      multiple
      clearable
      :items="items as string[]"
      :loading="pending as boolean"
      variant="outlined"
      hide-details
      label="Categories"
      closable-chips
      v-model="selected"
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
  </v-autocomplete>
</template>

<script lang="ts" setup>
const props = defineProps<{ source: string; modelValue: string[] }>();

const selected = ref<string[]>(props.modelValue);
const loading = ref<boolean>(true);

const { data: items, pending } = useAPI(
    (api) => api.categoryTree(undefined, { sources: [props.source] }).then(list => list.length ? list[0].categories.map(c => c.name) : []),
    {
      default: () => [],
      finally: () => {
        loading.value = false;
      },
    }
);

const emit = defineEmits<{ 'update:modelValue': [kinds: string[]] }>()

watch(selected, (current) => {
  emit('update:modelValue', current)
});
</script>
