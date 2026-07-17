<template>
  <v-card>
    <v-toolbar color="transparent">
      <v-toolbar-title>Cluster Resources</v-toolbar-title>
      <template #append>
        <Search v-model="search" style="min-width: 400px;" />
        <CollapseBtn v-model="open" />
      </template>
    </v-toolbar>
    <v-list v-if="pending" lines="two" class="mt-0 pt-0">
      <v-skeleton-loader class="mx-auto border" type="list-item-avatar" />
      <v-skeleton-loader class="mx-auto border" type="list-item-avatar" />
      <v-skeleton-loader class="mx-auto border" type="list-item-avatar" />
    </v-list>
    <template v-else>
      <v-list v-if="data?.items?.length && open" lines="two">
        <resource-item v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" :exceptions="exceptions" :show-skipped="showSkipped" />
      </v-list>
      <template v-if="data.count > options.offset">
        <v-divider />
        <v-pagination v-model="options.page" :length="length" class="my-4" />
      </template>
      <template v-if="!data?.items?.length">
        <v-divider />
        <v-card-text>
          No resources for the selected kinds found
        </v-card-text>
      </template>
    </template>
  </v-card>
</template>

<script setup lang="ts">
import type { Ref } from "vue";
import { Status, type Filter, type Pagination } from "~/types/core";
import { ClusterKinds, APIFilter } from "~/provider/dashboard";

const props = defineProps({
  details: { type: Boolean, required: true },
  exceptions: { type: Boolean, required: false },
  perPage: { type: Number, required: false, default: 8 },
  to: { type: [String, Object], required: false },
  id: { type: String, required: true },
});

const search = ref('')
const open = ref(true)

const filter = inject<Ref<Filter>>(APIFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(ClusterKinds, ref<string[]>([]))

const combinedFilter = computed(() => ({
  kinds: kinds.value.length ? kinds.value : undefined,
  search: search.value,
}))

const options = reactive<Pagination>({
  page: 1,
  offset: props.perPage,
})

const length = computed(() => Math.ceil((data.value?.count || 0) / options.offset))

const { data, refresh, pending } = useAPI(
    (api) => api.customBoardClusterResourceResults(props.id, combinedFilter.value, options),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

const status = useStatusInjection()
const showSkipped = computed(() => data.value.items.some(item => status.value?.includes(Status.SKIP) && !!item.status[Status.SKIP]))

watch(combinedFilter, onChange(() => {
  if (options.page !== 1) {
    options.page = 1
    return
  }

  refresh()
}))

watch(options, () => refresh())

</script>
