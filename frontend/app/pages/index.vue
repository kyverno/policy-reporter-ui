<template>
  <page-layout title="Dashboard"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               v-model:mode="mode"
               :source="source"
               v-if="data"
  >
    <template v-if="isCompact">
      <LazyGraphSourceCard v-if="data.singleSource" :data="data.charts.clusterScope[source!]!" :source="source!" :hide-cluster="!data.clusterScope" :type="dataType" />
      <LazyGraphSourcesStatusCard v-else :data="data" :hide-cluster="!data.clusterScope" />
      <resource-summary-list class="mt-6" :data="data.summary" />
    </template>
    <template v-else>
      <GraphSourceCharts :data="data" />
      <template v-for="source in data.showResults" :key="source">
        <app-row>
          <results :source="source" :title="`${capilize(source)}: Results without resource information`" />
        </app-row>
      </template>
      <app-row v-if="data.clusterScope">
        <resource-cluster-list :details="data.multiSource"  :exceptions="data.exceptions" />
      </app-row>
      <resource-namespace-section v-if="data.namespaces.length" :namespaces="data.namespaces">
        <template #default="{ namespaces }">
          <resource-scroller :list="namespaces">
            <template #default="{ item }">
              <resource-list :namespace="item" :details="data.multiSource" :exceptions="data.exceptions" />
            </template>
          </resource-scroller>
        </template>
      </resource-namespace-section>
    </template>
  </page-layout>
</template>

<script setup lang="ts">
import { APIFilter } from "~/provider/dashboard";

const router = useRouter()
const store = useSourceStore()
await store.load()

const { kinds, clusterKinds, filter } = useFilter()

const { data, refresh, error } = useAPI(api => api.dashboard(filter.value))

const { dataType, mode, isCompact } = useDashboardHelper(data)

watch(error, (err) => {
  if (err && err.status !== 401) {
    return;
  }

  callAPI((api) => api.layout()).then(layout => {
    if (!layout?.customBoards?.length) return;

    // @ts-ignore 
    router.push(layout.customBoards[0].path)
  })
})

watch(filter, onChange(refresh))

provide(APIFilter, filter)

const source = computed(() => {
  if (data.value?.sources.length !== 1) return undefined

  return data.value?.sources[0]
})

watch(source, (s?: string) => {
  if (!s) return
  
  store.load(s)
})

useSourceContext(source)
useStatusProvider(data)
useSeveritiesProvider(data)
</script>
