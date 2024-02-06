<template>
  <page-layout v-if="data"
               v-model:cluster-kinds="clusterKinds"
               v-model:kinds="kinds"
               :source="route.params.source"
               :category="route.params.category"
               :title="`${capilize(route.params.source)}: ${ route.params.category }`"
  >
    <GraphSourceStatus :category="route.params.category" :data="data" :source="route.params.source"/>
    <template v-if="data.showResults.length === 0">
      <app-row>
        <resource-cluster-list :details="data.multiSource" />
      </app-row>
      <resource-namespace-section v-if="data.namespaces.length" :namespaces="data.namespaces" />
    </template>
    <template v-else>
      <policy-cluster-results :source="route.params.source" :policy="route.params.policy" />
      <policy-namespace-section :namespaces="data.namespaces" :source="route.params.source" :policy="route.params.policy" />
    </template>
  </page-layout>
</template>

<script lang="ts" setup>
import { useAPI } from '~/modules/core/composables/api'
import { capilize } from "~/modules/core/layouthHelper";
import { APIFilter } from "~/modules/core/provider/dashboard";
import { onChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const route = useRoute()

const filter = computed(() => ({
  sources: [route.params.source],
  categories: [route.params.category],
  kinds: kinds.value,
  clusterKinds: clusterKinds.value,
}))

const { data, refresh } = useAPI((api) => api.dashboard(filter.value));

watch(filter, onChange(refresh))

provide(APIFilter, ref(filter))
</script>
