<template>
  <page-layout v-if="data"
               v-model:cluster-kinds="clusterKinds"
               v-model:kinds="kinds"
               :source="source"
               :category="category"
               :title="`${capilize(source)}: ${ route.params.category }`"
  >
    <GraphSourceSeverities v-if="data.type === 'severity'" :data="data" :source="source" :category="category" />
    <GraphSourceStatus v-else :data="data" :source="source" :category="category" />
    <template v-if="data.showResults.length === 0">
      <app-row v-if="data.clusterScope">
        <resource-cluster-list :details="false" :exceptions="data.exceptions" />
      </app-row>
      <resource-namespace-section v-if="data.namespaces.length" :namespaces="data.namespaces" :exceptions="data.exceptions" />
    </template>
    <template v-else>
      <policy-cluster-results :source="source" />
      <policy-namespace-section :namespaces="data.namespaces" :source="source" :exceptions="data.exceptions" />
    </template>
  </page-layout>
</template>

<script lang="ts" setup>

import { APIFilter } from "~/provider/dashboard";

import type {Filter} from "~/types/core";

const route = useRoute()

const source = computed(() => route.params.source as string)
const category = computed(() => route.params.category as string)

const store = useSourceStore(source.value)
await store.load(source.value)

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed((): Filter => ({
  sources: [source.value],
  categories: [category.value],
  kinds: kinds.value,
  clusterKinds: clusterKinds.value,
}))

const { data, refresh } = useAPI((api) => api.dashboard(filter.value));

watch(filter, onChange(refresh))

provide(APIFilter, filter)
useStatusProvider(data)
useSeveritiesProvider(data)
useSourceContext(source)
useDashboardType(data)
</script>
