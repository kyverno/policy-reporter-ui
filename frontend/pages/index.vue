<template>
  <v-container v-if="data.counts.length" fluid class="py-4 px-4 main-height">
    <v-row>
      <v-col>
        <v-toolbar color="indigo" elevation="2" rounded>
          <template #append>
            <FormKindAutocomplete style="min-width: 300px; max-width: 100%; margin-right: 15px;" v-model="kinds" />
            <FormClusterKindAutocomplete style="min-width: 300px;" v-model="clusterKinds" />
          </template>
        </v-toolbar>
      </v-col>
    </v-row>
    <SourcesStatus v-if="sources.length > 1" :data="data as FindingCounts" />
    <SourceStatus v-if="sources.length === 1" :data="data.counts[0]" />
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
  </v-container>
</template>

<script setup lang="ts">
import { callAPI, useAPI } from '~/modules/core/composables/api'
import { type FindingCounts } from "~/modules/core/types";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { ClusterKinds, NamespacedKinds } from "~/modules/core/provider/dashboard";
import { execOnChange } from "~/helper/compare";

const [sources, namespaces] = await Promise.all([
  callAPI((api) => api.sources().then((source) => source.map(s => s.name))),
  callAPI((api) => api.namespaces())
])

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ kinds: [...kinds.value, ...clusterKinds.value] }))

const { data, refresh } = useAPI((api) => api.countFindings(filter.value), {
  default: () => ({ total: 0, counts: [] })
});

watch(filter, (n, o) => execOnChange(n, o, refresh))

provide(NamespacedKinds, kinds)
provide(ClusterKinds, clusterKinds)
</script>
