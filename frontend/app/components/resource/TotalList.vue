<template>
  <v-card class="mb-4">
    <v-toolbar color="transparent">
      <v-toolbar-title>{{ cluster?.name ?? 'Namespace' }} Summary</v-toolbar-title>
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
        <resource-item v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" />
      </v-list>
      <template v-if="paginated">
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
import { type Filter, type Pagination } from '~/types/core'
import CollapseBtn from "~/components/CollapseBtn.vue";
import { APIFilter } from "~/provider/dashboard";
import type { Ref } from "vue";

const props = defineProps<{ cluster?: { name: string; slug: string }; details: boolean; perPage?: number }>()

const search = ref('')
const open = ref(true)
const filter = inject<Ref<Filter>>(APIFilter, ref<Filter>({}))

const options = reactive<Pagination>({
  page: 1,
  offset: props.perPage ?? 8,
})

const paginated = computed(() => {
  if (!data.value) return false

  if (options.offset < 1) return false

  return data.value.count > options.offset && open.value
})

const length = computed(() => {
  return Math.ceil((data.value?.count || 0) / options.offset)
})

const combinedFilter = computed(() => ({
  ...filter.value,
  search: search.value,
}))

const { data, refresh, pending } = useAPI(
    (api) => api.totalResults(props.cluster?.slug, options, combinedFilter.value),
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
  })
)

watch(options, () => refresh())

</script>
