<template>
  <div v-if="(data && data.count > 0) || !!searchText">
    <app-toolbar :title="category || `Results for ${capilize(source)}`" v-model:searchText="searchText" v-model:open="open" />
    <v-expand-transition>
      <div v-show="open">
        <v-divider />
        <resource-results-table
          :source="source"
          :category="category"
          :exceptions="exceptions"
          :resource="resource"
          :data="data"
          v-model:items-per-page="options.itemsPerPage"
          v-model:page="options.page"
          @search="searchText = $event"
        />
      </div>
    </v-expand-transition>
  </div>
</template>

<script lang="ts" setup>
import { Status } from "~/types/core";

const props = defineProps<{
  id: string;
  plugin?: boolean;
  source: string;
  category?: string;
  exceptions?: boolean;
  resource: string;
  status?: Status;
}>()

const options = reactive({
  itemsPerPage: 10,
  page: 1,
  sortDesc: [],
  sortBy: [],
  groupBy: [],
  groupDesc: [],
  multiSort: false,
  mustSort: false
})

const open = ref(true)
const searchText = ref('')
const status = useStatusInjection()

const { data, refresh } = useAPI(
    (api) => api.customBoardResourceDetailedResults(props.id, props.resource, {
      page: options.page,
      offset: options.itemsPerPage,
    }, {
      sources: props.source ? [props.source] : undefined,
      categories: props.category ? [props.category] : undefined,
      search: !!searchText.value ? searchText.value : undefined,
      status: props.status ? [props.status] : status.value,
    }),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(() => options.page, refresh)
watch(() => options.itemsPerPage, refresh)
watch(searchText, refresh)
</script>
