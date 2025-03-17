<template>
  <v-card v-if="(data && !pending) || !!searchText">
    <v-toolbar color="secondary">
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
              <td>{{ item.apiVersion }}</td>
              <td>{{ item.kind }}</td>
              <td>{{ item.name }}</td>
              <td>
                <v-btn v-if="plugin" @click.stop="openPolicy(item.policy, item.source)" class="mr-1" target="_blank" icon="mdi-open-in-new" variant="text" size="small" />
                {{ item.policy }}
              </td>
              <td>{{ item.rule }}</td>
              <td>
                <chip-severity v-if="item.severity" @click.prevent.stop="searchText = item.severity" :severity="item.severity" />
              </td>
              <td>
                <chip-status @click.prevent.stop="searchText = item.status" :status="item.status" />
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
  </v-card>
</template>

<script lang="ts" setup>
import { Status, type ResultList } from "~/modules/core/types";
import { mapResults } from "~/modules/core/mapper";
import { ClusterKinds } from "../../provider/dashboard";

const props = defineProps<{
  plugin?: boolean;
  sources: string[];
  exceptions?: boolean;
  Status?: Status;
}>()

const bg = useBGColor()

const router = useRouter()
const openPolicy = (policy: string, source: string) => {
  const { href } = router.resolve({ name: 'policies-source-info-policy', params: { source, policy }})
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

const kinds = inject<Ref<string[]>>(ClusterKinds, ref<string[]>([]))

const open = ref(true)
const searchText = ref('')
const { data, refresh, pending } = useAPI(
    (api) => api.clusterResults({
      sources: props.sources,
      kinds: kinds.value,
      search: !!searchText.value ? searchText.value : undefined,
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
watch(searchText, () => refresh())

const results = computed(() => mapResults(data.value as ResultList))

const headers = computed(() => {
  if (props.exceptions) {
    return [
    { title: 'APIVersion', key: 'apiVersion', width: '10%' },
    { title: 'Kind', key: 'kind', width: '10%' },
    { title: 'Name', key: 'name', width: '20%' },
    { title: 'Policy', key: 'policy', width: '18%' },
    { title: 'Rule', key: 'rule', width: '18%' },
    { title: 'Severity', key: 'severity', width: '9%' },
    { title: 'Status', key: 'status', width: '9%' },
    { title: 'Actions', key: 'exception', width: '9%', sortable: false }
    ]
  }

  return [
    { title: 'APIVersion', key: 'apiVersion', width: '10%' },
    { title: 'Kind', key: 'kind', width: '10%' },
    { title: 'Name', key: 'name', width: '20%' },
    { title: 'Policy', key: 'policy', width: '18%' },
    { title: 'Rule', key: 'rule', width: '18%' },
    { title: 'Severity', key: 'severity', width: '12%' },
    { title: 'Status', key: 'status', width: '12%' }
  ]
})
</script>
