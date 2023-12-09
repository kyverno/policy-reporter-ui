<template>
  <v-autocomplete
    multiple
    clearable
    density="compact"
    :items="items as string[]"
    variant="outlined"
    hide-details
    label="Kinds"
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
  </v-autocomplete>
</template>

<script lang="ts" setup>
import { kinds } from "~/modules/core/store/filter";

const props = defineProps({
  source: { type: String, default: undefined },
});

const selected = ref<string[]>(kinds.value);
const loading = ref<boolean>(true);

const { data: items } = useAPI(
  ($coreAPI) => {
    return $coreAPI.namespacedKinds(props.source)
  },
  {
    default: () => [],
    finally: () => {
      loading.value = false;
    },
  }
);

const input = defineRouteQuery('kinds', selected);

watch(items, (current) => {
  input(selected.value.filter((s) => current.includes(s)));
});

watch(selected, (current) => {
  kinds.value = current as string[]
});
</script>
