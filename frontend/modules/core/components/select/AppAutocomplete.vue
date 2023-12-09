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
import type { PropType } from "vue";
import type { NamespacedFilterAPI, ClusterFilterAPI } from "../../types";

const props = defineProps({
  namespaced: { type: Boolean, default: false },
  source: { type: String, default: undefined },
  label: { type: String, required: true },
  queryName: { type: String, required: true },
  namespacedApi: { type: String as PropType<NamespacedFilterAPI>, required: false },
  clusterApi: { type: String as PropType<ClusterFilterAPI>, required: false },
});

const emit = defineEmits(["itemsLoaded"]);

const selected = ref<string[]>([]);
const loading = ref<boolean>(true);

const { data: items, refresh } = useAPI(
  ($coreAPI) => {
    if (props.namespaced) {
      return $coreAPI[props.namespacedApi](props.source);
    }

    return $coreAPI[props.clusterApi](props.source);
  },
  {
    default: () => [],
    finally: () => {
      loading.value = false;
    },
  }
);

const input = defineRouteQuery(props.queryName, selected);

watch(items, (current) => {
  input(selected.value.filter((s) => current.includes(s)));
});

watch(selected, (current) => {
  emit("itemsLoaded", current);
});

doAPIRefresh(refresh);
</script>
