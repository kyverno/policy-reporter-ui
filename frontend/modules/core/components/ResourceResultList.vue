<template>
  <wait>
    <v-card>
      <v-toolbar color="transparent">
        <v-toolbar-title>{{ namespace }}</v-toolbar-title>
        <template #append>
          <Search class="mr-2" v-model="search" style="min-width: 300px;" />
          <CollapseBtn v-model="open" :disabled="!data.items.length" />
        </template>
      </v-toolbar>
      <v-list v-if="data?.items?.length && open" lines="two">
        <ResourceResultItem v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" />
      </v-list>
      <template v-if="!pending && !(data.items.length)">
        <v-divider />
        <v-card-text>
          No resources for the selected kinds found
        </v-card-text>
      </template>
    </v-card>
  </wait>
</template>

<script setup lang="ts">
import CollapseBtn from "~/components/CollapseBtn.vue";
import type { Filter } from "~/modules/core/types";
import type { Ref } from "vue";
import { NamespacedKinds, ResourceFilter } from "~/modules/core/provider/dashboard";
import { execOnChange } from "~/helper/compare";

const props = defineProps<{
  namespace: string;
  details: boolean;
}>()

const search = ref('')
const open = ref(true)

const filter = inject<Ref<Filter>>(ResourceFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(NamespacedKinds, ref<string[]>([]))

const combinedFilter = computed(() => ({
  ...filter.value,
  namespaces: [props.namespace as string],
  kinds: kinds.value.length ? kinds.value : undefined,
  search: search.value,
}))

const { data, refresh, pending } = useAPI(
    (api) => api.namespacedResourceResults(combinedFilter.value),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(combinedFilter, (n, o) => execOnChange(n, o, refresh))
</script>
