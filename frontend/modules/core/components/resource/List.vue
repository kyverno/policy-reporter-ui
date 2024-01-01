<template>
  <v-card>
    <v-toolbar color="transparent">
      <v-toolbar-title>{{ namespace }}</v-toolbar-title>
      <template #append>
        <Search class="mr-2" v-model="search" style="min-width: 300px;" />
        <CollapseBtn v-model="open" :disabled="!data.items.length" />
      </template>
    </v-toolbar>
    <v-list v-if="data?.items?.length && open" lines="two" class="pt-0">
      <resource-item v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" />
    </v-list>
    <template v-if="data.count > options.offset">
      <v-divider />
      <v-pagination v-model="options.page" :length="length" class="my-4" />
    </template>
    <template v-if="!pending && !(data.items.length)">
      <v-divider />
      <v-card-text>
        No resources for the selected kinds found
      </v-card-text>
    </template>
  </v-card>
</template>

<script setup lang="ts">
import type { Ref } from "vue";
import CollapseBtn from "~/components/CollapseBtn.vue";
import type { Filter, Pagination } from "~/modules/core/types";
import { NamespacedKinds, APIFilter } from "~/modules/core/provider/dashboard";
import { onChange } from "~/helper/compare";

const props = defineProps<{
  namespace: string;
  details: boolean;
}>()

const search = ref('')
const open = ref(true)

const filter = inject<Ref<Filter>>(APIFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(NamespacedKinds, ref<string[]>([]))

const combinedFilter = computed(() => ({
  ...filter.value,
  namespaces: [props.namespace as string],
  kinds: kinds.value.length ? kinds.value : undefined,
  search: search.value,
}))

const options = reactive<Pagination>({
  page: 1,
  offset: 8,
})

const length = computed(() => {
  return Math.ceil((data.value?.count || 0) / options.offset)
})

const { data, refresh, pending } = useAPI(
    (api) => api.namespacedResourceResults(combinedFilter.value, options),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(combinedFilter, onChange(() => {
  if (options.page !== 1) {
    options.page = 1
    return
  }

  refresh()
}))

watch(options, onChange<Pagination>(refresh))

</script>
