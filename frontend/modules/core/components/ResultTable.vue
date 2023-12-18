<template>
  <v-row v-if="data.count > 0 || !!searchText">
    <v-col cols="12">
      <v-card>
        <template #title>
          <span v-if="category">{{ category }}</span>
          <span v-else>Results for {{ capilize(source) }}</span>
        </template>
        <template #append>
          <search v-model="searchText" class="mr-4" style="min-width: 400px; float: left; height: 48px;" />
          <v-btn :icon="open ? 'mdi-chevron-up' : 'mdi-chevron-down'" @click="open = !open" variant="text" />
        </template>
        <v-expand-transition>
          <div v-show="open">
            <v-divider />
            <v-data-table-server
              :items="results.results"
              :items-length="results.count"
              :headers="headers"
              item-value="id"
              show-expand
              v-model:items-per-page="options.itemsPerPage"
              v-model:page="options.page"
            >
              <template #item.status="{ value }">
                <chip-status @click="searchText = value" :status="value" />
              </template>
              <template #item.severity="{ value }">
                <chip-severity v-if="value" @click="searchText = value" :severity="value" />
              </template>
              <template #expanded-row="{ columns, item }">
                <tr class="table-expand-text">
                  <td :colspan="columns.length" class="py-3">
                    <div v-if="item.hasProps">
                      <v-card v-if="item.message" variant="flat">
                        <v-alert type="info" variant="flat">
                          {{ item.message }}
                        </v-alert>
                      </v-card>
                      <div>
                        <template v-for="(value, label) in item.chips"  :key="label">
                          <property-chip :label="label as string" :value="value" class="mr-2 mt-2 rounded-xl" />
                        </template>
                        <template v-for="(value, label) in item.cards"  :key="label">
                          <property-card :label="label as string" :value="value" class="mt-2" />
                        </template>
                      </div>
                    </div>
                    <div v-else>
                      {{ item.message }}
                    </div>
                  </td>
                </tr>
              </template>
            </v-data-table-server>
          </div>
        </v-expand-transition>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
import { type Dictionary, Status } from "~/modules/core/types";
import { capilize } from "~/modules/core/layouthHelper";
import { mapResults } from "~/modules/core/mapper";

const props = defineProps<{
  source: string;
  category?: string;
  resource: string;
  Status?: Status;
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

const { data, refresh } = useAPI(
    (api) => api.results(props.resource, {
      page: options.page,
      offset: options.itemsPerPage,
    }, {
      sources: props.source ? [props.source] : undefined,
      categories: props.category ? [props.category] : undefined,
      search: searchText.value,
    }),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(() => options.page, refresh)
watch(() => options.itemsPerPage, refresh)
watch(searchText, refresh)

const results = computed(() => mapResults(data.value))

const headers = [
  { title: 'Policy', key: 'policy', width: '33%' },
  { title: 'Rule', key: 'rule', width: '33%' },
  { title: 'Severity', key: 'severity', width: '17%' },
  { title: 'Status', key: 'status', width: '17%' }
]
</script>

<style>
.table-expand-text {
  background-color: #f5f5f5;

  &:hover {
    background-color: #f5f5f5 !important;
  }
}
</style>
