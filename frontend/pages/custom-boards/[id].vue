<template>
  <page-layout v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               v-if="data"
  >
    <GraphSourceStatus v-if="data.singleSource" :data="data" :source="data.sources[0]" />
    <GraphSourcesStatus v-else :data="data" :hide-cluster="!data.clusterScope" />
    <v-row v-if="data.clusterScope">
      <v-col>
        <resource-cluster-list :details="data.multiSource" />
      </v-col>
    </v-row>
    <resource-namespace-section v-if="data.namespaces.length" :namespaces="data.namespaces">
      <template #default="{ namespaces }">
        <resource-scroller :list="namespaces" :default-loadings="3">
          <template #default="{ item }">
            <resource-list :namespace="item" :details="data.multiSource" />
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

const route = useRoute()

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({
  kinds: kinds.value,
  clusterKinds: clusterKinds.value
}))

const { data, refresh } = useAPI((api) => api.customBoard(route.params.id, filter.value))

watch(filter, onChange(refresh))

provide(APIFilter, computed(() => ({
  ...filter.value,
  sources: data.value?.filterSources,
})))
</script>
