<template>
  <page-layout :title="capilize(route.params.source)" v-model:kinds="kinds" v-model:cluster-kinds="clusterKinds" :source="route.params.source" v-if="data">
    <GraphSourceStatus :data="data" :source="route.params.source" />
    <template v-if="data.showResults.length === 0">
      <v-row>
        <v-col>
          <resource-cluster-list :source="route.params.source" :details="false" />
        </v-col>
      </v-row>
      <resource-namespace-section v-if="data.namespaces" :namespaces="data.namespaces" />
    </template>
    <template v-else>
      <policy-cluster-results :source="route.params.source" :policy="route.params.policy" />
      <policy-namespace-section :namespaces="data.namespaces" :source="route.params.source" :policy="route.params.policy" />
    </template>
  </page-layout>
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import { type Filter } from "~/modules/core/types";
import { APIFilter } from "~/modules/core/provider/dashboard";
import { onChange } from "~/helper/compare";

const route = useRoute()

const store = useSourceStore(route.params.source)
await store.load()

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed<Filter>(() => ({
  sources: [route.params.source],
  kinds: kinds.value,
  clusterKinds: clusterKinds.value,
}))

const { data, refresh } = useAPI((api) => api.dashboard(filter.value))

watch(filter, onChange(refresh))

provide(APIFilter, filter)
</script>
