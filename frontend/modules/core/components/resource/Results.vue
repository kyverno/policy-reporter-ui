<template>
  <div v-if="(data && data.count > 0) || !!searchText">
    <v-toolbar color="secondary">
      <template #title>
        <span v-if="category">{{ category }}</span>
        <span v-else>Results for {{ capilize(source) }}</span>
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
          hover
        >
          <template #item="{ item, ...props }">
            <tr @click="() => props.toggleExpand(props.internalItem)" class="cursor-pointer">
              <td>
                <v-btn v-if="plugin" @click.stop="openPolicy(item.policy)" class="mr-1" target="_blank" icon="mdi-open-in-new" variant="text" size="small" />
                {{ item.policy }}
              </td>
              <td>{{ item.rule }}</td>
              <td>
                <chip-severity v-if="item.severity" @click.prevent.stop="searchText = item.severity" :severity="item.severity" />
              </td>
              <td>
                <chip-status @click.prevent.stop="searchText = item.status" :status="item.status" />
              </td>
              <td>
                <resource-exception-dialog v-if="exceptions" :resource="resource" :source="source" :policies="[{ name: item.policy, rules: [{ name: item.rule, props: item.properties }]}]" />
              </td>
            </tr>
          </template>
          <template #expanded-row="{ columns, item }">
            <tr :class="bg">
              <td :colspan="columns.length" class="py-3">
                <div v-if="item.hasProps">
                  <v-card v-if="item.message" variant="flat">
                    <v-alert type="info" variant="flat" class="text-pre-line">
                      {{ item.message }}
                    </v-alert>
                  </v-card>
                  <div>
                    <template v-for="(value, label) in item.chips"  :key="label">
                      <property-chip :label="label as string" :value="value" class="mr-2 mt-2 rounded-lg" />
                    </template>
                    <template v-for="(value, label) in item.cards"  :key="label">
                      <property-card :label="label as string" :value="value" class="mt-2" />
                    </template>
                  </div>
                </div>
                <div class="text-pre-line" v-else>
                  {{ item.message }}
                </div>
              </td>
            </tr>
          </template>
        </v-data-table-server>
      </div>
    </v-expand-transition>
  </div>
</template>

<script lang="ts" setup>
import { Status, type ResultList } from "~/modules/core/types";
import { capilize } from "~/modules/core/layouthHelper";
import { mapResults } from "~/modules/core/mapper";

const props = defineProps<{
  plugin?: boolean;
  source: string;
  category?: string;
  exceptions?: boolean;
  resource: string;
  Status?: Status;
}>()

const bg = useBGColor()

const router = useRouter()
const openPolicy = (policy: string) => {
  const { href } = router.resolve({ name: 'policies-source-info-policy', params: { source: props.source, policy }})
  window.open(href, '_blank')
}


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
      search: !!searchText.value ? searchText.value : undefined,
    }),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(() => options.page, () => refresh())
watch(() => options.itemsPerPage, () => refresh())
watch(searchText, () => refresh())

const results = computed(() => mapResults(data.value as ResultList))

const headers = computed(() => {
  if (props.exceptions) {
    return [
      { title: 'Policy', key: 'policy', width: '33%' },
      { title: 'Rule', key: 'rule', width: '33%' },
      { title: 'Severity', key: 'severity', width: '12%' },
      { title: 'Status', key: 'status', width: '12%' },
      { title: 'Actions', key: 'exception', width: '12%', sortable: false }
    ]
  }

  return [
    { title: 'Policy', key: 'policy', width: '33%' },
    { title: 'Rule', key: 'rule', width: '33%' },
    { title: 'Severity', key: 'severity', width: '17%' },
    { title: 'Status', key: 'status', width: '17%' }
  ]
})
</script>
