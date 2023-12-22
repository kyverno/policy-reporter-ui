<template>
  <v-container fluid v-if="data.counts.length" class="py-4 px-4 main-height">
    <v-row>
      <v-col>
        <v-toolbar color="indigo" elevation="2" rounded>
          <template #append>
            <FormKindAutocomplete v-model="kinds" style="min-width: 300px; max-width: 100%; margin-left: 15px;" />
          </template>
        </v-toolbar>
      </v-col>
    </v-row>
    <SourcesStatus v-if="sources.length > 1" :hide-cluster="true" :data="data as FindingCounts" />
    <SourceStatus v-if="sources.length === 1" :hide-cluster="true" :data="data.counts[0]" />
    <resource-scroller v-if="details" :list="namespaces">
      <template #default="{ item }">
        <LazyResourceResultList :namespace="item" :details="sources.length > 1" :filter="filter" />
      </template>
    </resource-scroller>
  </v-container>
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

provide(ResourceFilter, filter)
provide(NamespacedKinds, kinds)
</script>
