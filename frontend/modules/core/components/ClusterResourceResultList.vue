<template>
<v-card>
  <v-toolbar color="transparent">
    <v-toolbar-title>Cluster Resources</v-toolbar-title>
    <template #append>
      <Search v-model="search" style="min-width: 300px;" />
      <CollapseBtn v-model="open" />
    </template>
  </v-toolbar>
  <v-list v-if="data?.items?.length && open" lines="two">
    <ResourceResultItem v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" />
  </v-list>
  <template v-if="data.count > options.offset">
    <v-divider />
    <v-pagination v-model="options.page" :length="length" class="my-4" />
  </template>
  <template v-if="!pending && !(data?.items?.length)">
      <v-divider />
      <v-card-text>
          No resources for the selected kinds found
      </v-card-text>
  </template>
</v-card>
</template>

<script setup lang="ts">
import { type Filter, type Pagination } from '../types'
import CollapseBtn from "~/components/CollapseBtn.vue";
import type { Ref } from "vue";
import { ClusterKinds, ResourceFilter } from "~/modules/core/provider/dashboard";
import { execOnChange } from "~/helper/compare";

const props = defineProps<{ source?: string; details: boolean; }>()

const search = ref('')
const open = ref(true)
const filter = inject<Ref<Filter>>(ResourceFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(ClusterKinds, ref<string[]>([]))

const options = reactive<Pagination>({
  page: 1,
  offset: 8,
})

const length = computed(() => {
  return Math.ceil((data.value?.count || 0) / options.offset)
})

const combinedFilter = computed(() => ({
  ...filter.value,
  kinds: kinds.value.length ? kinds.value : undefined,
  search: search.value,
}))

const { data, refresh, pending } = useAPI(
    (api) => api.clusterResourceResults(combinedFilter.value, options),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(combinedFilter, (o, n) =>
  execOnChange(o, n, () => {
    if (options.page !== 1) {
      options.page = 1
      return
    }

    refresh()
  })
)

watch(options, () => refresh())

</script>
