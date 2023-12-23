<template>
  <page-layout v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
  >
    <SourcesStatus v-if="sources.length > 1" :data="data as FindingCounts" />
    <SourceStatus v-if="sources.length === 1" :data="findings" />
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :details="sources.length > 1" />
      </v-col>
    </v-row>
    <resource-scroller :list="namespaces" v-if="namespaces.length">
      <template #default="{ item }">
        <LazyResourceResultList :namespace="item" :details="sources.length > 1" />
      </template>
    </resource-scroller>
  </page-layout>
</template>

<script setup lang="ts">
import { callAPI, useAPI } from '~/modules/core/composables/api'
import { type FindingCounts } from "~/modules/core/types";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { execOnChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const [sources, namespaces] = await Promise.all([
  callAPI((api) => api.sources().then((source) => source.map(s => s.name))),
  callAPI((api) => api.namespaces())
])

const filter = computed(() => ({ kinds: [...kinds.value, ...clusterKinds.value] }))

const { data, refresh } = useAPI((api) => api.countFindings(filter.value), {
  default: () => ({ total: 0, counts: [] })
});

const findings = computed(() => {
  if (data.value?.counts?.length) return data.value.counts[0];

  return { source: sources[0], counts: {}, total: 0 }
})

watch(filter, (n, o) => execOnChange(n, o, refresh))

</script>
