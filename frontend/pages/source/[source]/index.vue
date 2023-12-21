<template>
  <v-container fluid class="py-4 px-4 main-height">
    <v-row>
      <v-col>
        <v-toolbar color="indigo" elevation="2" rounded>
          <v-toolbar-title>{{ capilize(route.params.source) }}</v-toolbar-title>
          <template #append>
            <SelectKindAutocomplete style="width: 500px; max-width: 100%; margin-left: 15px;" :source="route.params.source" />
          </template>
        </v-toolbar>
      </v-col>
    </v-row>
    <SourceStatus v-if="data.counts.length" :data="data.counts[0]" />
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :filter="filter" :source="route.params.source" :details="false" />
      </v-col>
    </v-row>
    <v-infinite-scroll :onLoad="load" class="no-scrollbar" v-if="!pending && loaded.length">
      <template v-for="ns in loaded" :key="ns">
        <v-row>
          <v-col>
            <LazyResourceResultList :namespace="ns" :details="false" :filter="filter" />
          </v-col>
        </v-row>
      </template>
      <template #empty></template>
    </v-infinite-scroll>
  </v-container>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { clusterKinds, kinds } from '~/modules/core/store/filter';
import { capilize } from "~/modules/core/layouthHelper";
import { useInfinite } from "~/composables/infinite";
import type { Filter } from "~/modules/core/types";

const route = useRoute()

const { data, refresh } = useAPI(
    (api) => api.countFindings({ kinds: [...kinds.value, ...clusterKinds.value], sources: [route.params.source] }),
    {
      default: () => ({ total: 0, counts: [] }),
    }
);

const { data: namespaces, pending } = useAPI(
    (api) => api.namespaces(route.params.source),
    {
      default: () => [],
    }
);

const filter = computed<Filter>(() => ({
  sources: [route.params.source]
}))

const { load, loaded } = useInfinite(namespaces)

watch(kinds, refresh)
</script>
