<template>
  <v-container fluid v-if="data" class="py-4 px-4 main-height">
    <v-row>
      <v-col>
        <v-toolbar color="indigo" elevation="2" rounded>
          <template #append>
            <SelectKindAutocomplete style="width: 500px; max-width: 100%; margin-left: 15px;" />
          </template>
        </v-toolbar>
      </v-col>
    </v-row>
    <SourcesStatus v-if="multiSource" :hide-cluster="true" :data="data as FindingCounts" :filter="filter" />
    <SourceStatus v-if="singleSource && data.counts.length > 0" :hide-cluster="true" :data="data.counts[0]" :filter="filter" />
    <resource-scroller v-if="details" :details="multiSource" :list="details.namespaces" :filter="filter" />
  </v-container>
</template>

<script setup lang="ts">
import { callAPI, useAPI } from '~/modules/core/composables/api'
import { useInfinite } from "~/composables/infinite";
import { clusterKinds, kinds } from '~/modules/core/store/filter';
import { type FindingCounts } from "~/modules/core/types";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";

const route = useRoute()

const { data: details } = useAPI((api) => api.customBoard(route.params.id), { default: () => ({
    id: route.params.id,
    sources: [] as string[],
    namespaces: [] as string[],
  })
});

const { data: allSources } = useAPI((api) => api.sources().then((source) => source.map(s => s.name)), {
  default: () => []
});

const sources = computed<string[]>(() => {
  if (!details.value?.sources || !details.value?.sources.length) return allSources.value || []

  return allSources.value?.filter(s => details.value?.sources.some(d => s.toLowerCase() === d.toLowerCase())) || []
})

const multiSource = computed(() => (sources.value.length || 0) > 1)
const singleSource = computed(() => (sources.value.length || 0) === 1)

const namespaces = computed(() => details.value?.namespaces || [])

const filter = computed(() => ({
  kinds: [...kinds.value, ...clusterKinds.value],
  namespaces: details.value?.namespaces,
  sources: sources.value.length === allSources.value?.length ? undefined : sources.value
}))

const { data, refresh } = useAPI((api) => api.countFindings(filter.value), { default: () => ({ total: 0, counts: [] }) });

watch(kinds, refresh)
</script>
