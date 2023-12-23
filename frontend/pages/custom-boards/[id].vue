<template>
  <page-layout v-model:kinds="kinds"
               :title="details.name"
               :source="sources.length > 1 ? undefined : sources[0]"
               ns-scoped
  >
    <SourcesStatus v-if="sources.length > 1" :hide-cluster="true" :data="data as FindingCounts" />
    <SourceStatus v-if="sources.length === 1" :hide-cluster="true" :data="findings" />
    <resource-scroller v-if="details" :list="namespaces">
      <template #default="{ item }">
        <LazyResourceResultList :namespace="item" :details="sources.length > 1" :filter="filter" />
      </template>
    </resource-scroller>
  </page-layout>
</template>

<script setup lang="ts">
import { callAPI, useAPI } from '~/modules/core/composables/api'
import { type Filter, type FindingCounts } from "~/modules/core/types";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { NamespacedKinds, ResourceFilter } from "~/modules/core/provider/dashboard";

const route = useRoute()
const kinds = ref<string[]>([])

const loading = ref(true)

const { details, sources, namespaces, filterSource } = await Promise.all([
  callAPI((api) => api.customBoard(route.params.id)),
  callAPI((api) => api.sources().then((source) => source.map(s => s.name))),
  callAPI((api) => api.namespaces()),
]).then(([details, allSources, allNamespaces]) => {
  const namespaces = details?.namespaces || allNamespaces || []

  if (!details?.sources || !details?.sources.length) {
    return { details, sources: allSources, namespaces, filterSource: false }
  }

  const sources = allSources.filter(s => details?.sources.some(d => s.toLowerCase() === d.toLowerCase())) || []

  return {
    details,
    sources,
    namespaces,
    filterSource: allSources.length !== sources.length
  }
}).finally(() => { loading.value = false })


const filter = computed<Filter>(() => ({
  kinds: [...kinds.value],
  namespaced: true,
  namespaces: details?.namespaces?.length ? details?.namespaces : undefined,
  sources: filterSource ? sources : undefined,
  labels: Object.entries((details?.labels || {})).map(([label, value]) => `${label}:${value}`)
}))

const { data, refresh } = useAPI((api) => api.countFindings(filter.value), { default: () => ({ total: 0, counts: [] }) });

watch(filter, () => refresh())

const findings = computed(() => {
  if (data.value?.counts?.length) return data.value.counts[0];

  return { source: sources[0], counts: {}, total: 0 }
})

provide(ResourceFilter, filter)
</script>
