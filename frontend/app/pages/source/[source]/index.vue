<template>
  <page-layout
    :title="capilize(source)"
    v-model:kinds="kinds"
    v-model:cluster-kinds="clusterKinds"
    v-model:mode="mode"
    :source="source"
    v-if="data"
  >
    <template v-if="isCompact">
      <LazyGraphSourceCard
        :data="data.charts.clusterScope[source!]!"
        :source="source!"
        :hide-cluster="!data.clusterScope"
        :type="dataType"
      />
      <resource-summary-list class="mt-6" :data="data.summary" />
    </template>
    <template v-else>
      <LazyGraphSourceSeverities
        v-if="isSeverity"
        :data="data"
        :source="source"
      />
      <LazyGraphSourceStatus v-else :data="data" :source="source" />
      <template v-if="data.showResults.length === 0">
        <v-row v-if="data.clusterScope">
          <v-col>
            <resource-cluster-list
              :source="route.params.source"
              :details="false"
              :exceptions="data.exceptions"
            />
          </v-col>
        </v-row>
        <resource-namespace-section
          v-if="data.namespaces"
          :namespaces="data.namespaces"
          :exceptions="data.exceptions"
        />
      </template>
      <template v-else>
        <policy-cluster-results :source="source" />
        <policy-namespace-section
          :namespaces="data.namespaces"
          :source="source"
        />
      </template>
    </template>
  </page-layout>
</template>

<script setup lang="ts">
import { APIFilter } from '~/provider/dashboard';

const route = useRoute();

const { kinds, clusterKinds, filter, source } = useFilter()

const store = useSourceStore(source.value);
await store.load(source.value);


const { data, refresh } = useAPI((api) => api.dashboard(filter.value));
const { isSeverity, dataType, mode, isCompact } = useDashboardHelper(data)

watch(filter, onChange(refresh));

provide(APIFilter, filter);
useStatusProvider(data);
useSeveritiesProvider(data);
useSourceContext(source);
useDashboardType(data);
</script>
