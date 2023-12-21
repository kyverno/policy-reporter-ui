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
    <SourcesStatus v-if="multiSource" :data="data as FindingCounts" />
    <SourceStatus v-if="singleSource && data.counts.length > 0" :data="data.counts[0]" />
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :details="multiSource" />
      </v-col>
    </v-row>
    <resource-scroller :details="multiSource" :list="namespaces" v-if="namespaces.length" />
  </v-container>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { useInfinite } from "~/composables/infinite";
import { clusterKinds, kinds } from '~/modules/core/store/filter';
import { type FindingCounts } from "~/modules/core/types";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";

const { data: sources } = useAPI<string[]>((api) => api.sources().then((source) => source.map(s => s.name)), {
  default: () => []
});

const multiSource = computed(() => (sources.value?.length || 0) > 1)
const singleSource = computed(() => (sources.value?.length || 0) === 1)

const { data, refresh } = useAPI((api) => api.countFindings({ kinds: [...kinds.value, ...clusterKinds.value] }), {
  default: () => ({ total: 0, counts: [] }),
});

const { data: namespaces } = useAPI((api) => api.namespaces(), { default: () => [] });

watch(kinds, refresh)
</script>
