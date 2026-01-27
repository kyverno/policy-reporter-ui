<template>
  <page-layout title="Dashboard"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="data.singleSource ? data.sources[0] : undefined"
               v-if="data"
  >
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
  </page-layout>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { onChange } from "~/helper/compare";
import { APIFilter } from "~/modules/core/provider/dashboard";
import { capilize } from "~/modules/core/layouthHelper";
import { useConfigStore } from "~/store/config";

const router = useRouter()
const config = useConfigStore()

const store = useSourceStore()
await store.load()

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({
  kinds: kinds.value,
  clusterKinds: clusterKinds.value
}))

const { data, refresh, error } = useAPI(api => api.dashboard(filter.value))

watch(error, (err) => {
  if (err && err.status !== 401) {
    return;
  }

  callAPI((api) => api.layout()).then(layout => {
    if (!layout.customBoards.length) return;

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
