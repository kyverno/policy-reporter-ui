<template>
  <page-layout
    title="Dashboard"
    v-model:kinds="kinds"
    v-model:cluster-kinds="clusterKinds"
    :source="data.singleSource ? data.sources[0] : undefined"
    v-if="data"
  >
    <v-row v-if="data && data.charts.clusters">
      <v-col cols="12">
        <v-card style="height: 100%">
          <v-card-text style="height: 100%" v-if="showExpanded">
            <GraphCountPerNamespace
              :data="data.charts.clusters.complete"
              title="Results per Cluster"
            />
          </v-card-text>
          <v-card-text style="height: 100%" v-else>
            <GraphCountPerNamespace
              :data="data.charts.clusters.preview!"
              title="Results per Cluster"
            />
          </v-card-text>
          <v-btn
            v-if="hasPreview"
            variant="outlined"
            size="small"
            @click="expand = !expand"
            style="position: absolute; bottom: 10px; right: 10px"
            rounded="0"
          >
            <span v-if="showExpanded">Show preview</span>
            <span v-else>Show Complete List</span>
          </v-btn>
        </v-card>
      </v-col>
    </v-row>

    <resource-cluster-section
      v-if="config.clusters.length"
      :clusters="config.clusters"
    >
      <template #default="{ clusters }">
        <resource-scroller :list="clusters">
          <template #default="{ item }">
            <resource-total-list :cluster="item" :details="false" />
          </template>
        </resource-scroller>
      </template>
    </resource-cluster-section>
  </page-layout>
</template>

<script setup lang="ts">
import { APIFilter } from '~/provider/dashboard';
import { useConfigStore } from '~/store/config';
import type { Dashboard } from '~/types/core';

const config = useConfigStore();

const kinds = ref<string[]>([]);
const clusterKinds = ref<string[]>([]);

const filter = computed(() => ({
  kinds: kinds.value,
  clusterKinds: clusterKinds.value,
}));

const { data, refresh, error } = useAPI<Dashboard, null>((api) =>
  api.clustersDashboard(filter.value),
);

watch(filter, onChange(refresh));

provide(APIFilter, filter);

const expand = ref(false);
const hasPreview = computed(() => !!data.value?.charts.clusters?.preview);

const showExpanded = computed(() => {
  if (!data.value?.charts.clusters?.preview) {
    return true;
  }

  return expand.value;
});

const source = computed(() => {
  if (data.value?.sources.length !== 1) return undefined;

  return data.value?.sources[0];
});

useSourceContext(source);
useStatusProvider(data);
useSeveritiesProvider(data);
</script>
