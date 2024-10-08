<template>
  <app-row v-if="data.count > 0 || !!searchText">
    <v-card>
      <v-toolbar color="transparent">
        <template #title>
          <span>{{ namespace }}</span>
        </template>
        <template #append>
          <search v-model="searchText" class="mr-4" style="min-width: 400px; float: left;" />
          <v-btn :icon="open ? 'mdi-chevron-up' : 'mdi-chevron-down'" @click="open = !open" variant="text" />
        </template>
      </v-toolbar>
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
            <template #item.exception="{ item }" v-if="props.exceptions && props.policy">
              <resource-exception-dialog :resource="item.resourceId" :source="props.source" :policies="[{ name: props.policy, rules: [{ name: item.rule, props: item.properties }]}]" />
            </template>
            <template #expanded-row="{ columns, item }">
              <tr :class="bg">
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
            <template #bottom v-if="results.count <= options.itemsPerPage"></template>
          </v-data-table-server>
        </div>
      </v-expand-transition>
    </v-card>
  </app-row>
</template>

<script lang="ts" setup>
import { mapResults } from "~/modules/core/mapper";
import type { Ref } from "vue";
import { type Filter, Status } from "~/modules/core/types";
import { APIFilter } from "~/modules/core/provider/dashboard";
import { onChange } from "~/helper/compare";

const props = defineProps<{
  source: string;
  namespace: string;
  status?: Status[];
  policy?: string;
  exceptions?: boolean;
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
const bg = useBGColor()

const filter = inject<Ref<Filter>>(APIFilter, ref<Filter>({}))

const { data, refresh } = useAPI(
    (api) => api.namespacedResults({
      ...filter.value,
      sources: [props.source],
      policies: props.policy ? [props.policy] : undefined,
      status: props.status ? props.status : undefined,
      namespaces: [props.namespace],
      search: searchText.value ? searchText.value : undefined,
    }, {
      page: options.page,
      offset: options.itemsPerPage,
    }),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(() => options.page, () => refresh())
watch(() => options.itemsPerPage, () => refresh())
watch(() => props.status, () => refresh())
watch(searchText, () => refresh())

watch(filter, onChange(refresh))

const results = computed(() => mapResults(data.value))

const headers = computed(() => {
  const noResources = results.value.results.every(r => !r.name)

  if (noResources && !props.exceptions) {
    return [
      { title: 'Policy', key: 'policy', width: '36%' },
      { title: 'Rule', key: 'rule', width: '36%' },
      { title: 'Severity', key: 'severity', width: '12%' },
      { title: 'Status', key: 'status', width: '12%' },
    ]
  }

  if (noResources && props.exceptions && props.policy) {
    return [
      { title: 'Policy', key: 'policy', width: '30%' },
      { title: 'Rule', key: 'rule', width: '30%' },
      { title: 'Severity', key: 'severity', width: '12%' },
      { title: 'Status', key: 'status', width: '12%' },
      { title: 'Actions', key: 'exception', width: '12%' },
    ]
  }

  if (props.exceptions && props.policy) {
    return [
      { title: 'APIVersion', key: 'apiVersion', width: '12%' },
      { title: 'Kind', key: 'kind', width: '14%' },
      { title: 'Name', key: 'name', width: '22%' },
      { title: 'Rule', key: 'rule', width: '22%' },
      { title: 'Severity', key: 'severity', width: '9%' },
      { title: 'Status', key: 'status', width: '9%' },
      { title: 'Actions', key: 'exception', width: '12%' },
    ]
  }

  return [
    {title: 'APIVersion', key: 'apiVersion', width: '10%'},
    {title: 'Kind', key: 'kind', width: '16%'},
    {title: 'Name', key: 'name', width: '25%'},
    {title: 'Rule', key: 'rule', width: '25%'},
    {title: 'Severity', key: 'severity', width: '12%'},
    {title: 'Status', key: 'status', width: '12%'}
  ]
})
</script>

<style>
.table-expand-text {
  background-color: #f5f5f5;

  &:hover {
    background-color: #f5f5f5 !important;
  }
}
</style>
