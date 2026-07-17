<template>
  <page-layout v-if="data"
               v-model:cluster-kinds="clusterKinds"
               v-model:kinds="kinds"
               v-model:mode="mode"
               :source="source"
               :category="category"
               :title="`${capilize(source)}: ${ category }`"
  >
    <template v-if="isCompact">
      <LazyGraphSourceCard :data="data.charts.clusterScope[source!]!" :source="source!" :hide-cluster="!data.clusterScope" :type="dataType" />
      <resource-summary-list class="mt-6" :data="data.summary" :source="source" :category="category" />
    </template>
    <template v-else>
      <LazyGraphSourceSeverities v-if="isSeverity" :data="data" :source="source" :category="category" />
      <LazyGraphSourceStatus v-else :data="data" :source="source" :category="category" />
      <template v-if="data.showResults.length === 0">
        <app-row v-if="data.clusterScope">
          <resource-cluster-list :details="false" :exceptions="data.exceptions" />
        </app-row>
        <resource-namespace-section v-if="data.namespaces.length" :namespaces="data.namespaces" :exceptions="data.exceptions" />
      </template>
      <template v-else>
        <policy-cluster-results :source="source" />
        <policy-namespace-section :namespaces="data.namespaces" :source="source" :exceptions="data.exceptions" />
      </template>
    </template>
  </page-layout>
</template>

<script lang="ts" setup>

import { APIFilter } from "~/provider/dashboard";

const { kinds, clusterKinds, filter, source, category } = useFilter()

const { data, refresh } = useAPI((api) => api.dashboard(filter.value));
const { isSeverity, dataType, mode, isCompact } = useDashboardHelper(data)

watch(filter, onChange(refresh))

provide(APIFilter, filter)
useDashboardProvider(data)
</script>
