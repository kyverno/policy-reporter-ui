<template>
  <v-card v-if="data">
    <v-toolbar color="transparent">
      <v-toolbar-title>{{ namespace }}</v-toolbar-title>
      <template #append>
        <Search class="mr-2" v-model="search" style="min-width: 300px;" />
        <CollapseBtn v-model="open" :disabled="!(data?.items?.length)" />
      </template>
    </v-toolbar>
    <template v-if="open">
      <v-list v-if="data?.items?.length" lines="two" class="pt-0">
        <resource-item v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" :exceptions="exceptions" :show-skipped="showSkipped" />
      </v-list>
      <template v-if="data.count > options.offset">
        <v-divider />
        <v-pagination v-model="options.page" :length="length" class="my-4" />
      </template>
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
import { Status, type Filter, type Pagination } from "~/modules/core/types";
import { NamespacedKinds, APIFilter } from "~/modules/core/provider/dashboard";
import { onChange } from "~/helper/compare";

const props = defineProps<{
  namespace: string;
  details: boolean;
  exceptions?: boolean;
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

const length = computed(() => Math.ceil((data.value?.count || 0) / options.offset))

const { data, refresh, pending } = useAPI(
    (api) => api.namespacedResourceResults(combinedFilter.value, options),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

const status = useStatusInjection()
const showSkipped = computed(() => data.value?.items.some(item => status.value.includes(Status.SKIP) && !!item[Status.SKIP]))

watch(combinedFilter, onChange(() => {
  if (options.page !== 1) {
    options.page = 1
    return
  }

  refresh()
}))

watch(options, () => refresh())

</script>
