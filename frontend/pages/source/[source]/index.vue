<template>
  <page-layout :title="capilize(source)" v-model:kinds="kinds" v-model:cluster-kinds="clusterKinds" :source="source" v-if="data">
    <GraphSourceSeverities v-if="data.type === 'severity'" :data="data" :source="source" />
    <GraphSourceStatus v-else :data="data" :source="source" />
    <template v-if="data.showResults.length === 0">
      <v-row>
        <v-col>
          <resource-cluster-list :source="route.params.source" :details="false" :exceptions="data.exceptions" />
        </v-col>
      </v-row>
      <resource-namespace-section v-if="data.namespaces" :namespaces="data.namespaces" :exceptions="data.exceptions" />
    </template>
    <template v-else>
      <policy-cluster-results :source="source" />
      <policy-namespace-section :namespaces="data.namespaces" :source="source" />
    </template>
  </page-layout>
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import { type Filter } from "~/modules/core/types";
import { APIFilter } from "~/modules/core/provider/dashboard";
import { onChange } from "~/helper/compare";

const route = useRoute()

const source = computed(() => route.params.source as string)

const store = useSourceStore(source.value)
await store.load(source.value)

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed((): Filter => ({
  sources: [source.value],
  kinds: kinds.value,
  clusterKinds: clusterKinds.value,
}))

const { data, refresh } = useAPI((api) => api.dashboard(filter.value))

watch(filter, onChange(refresh))

provide(APIFilter, filter)
useStatusProvider(data)
useSeveritiesProvider(data)
useSourceContext(source)
useDashboardType(data)
</script>
