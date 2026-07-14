<template>
  <app-row v-if="!hideCluster">
    <v-card :title="`${capilize(source)} Cluster Scoped Results`">
      <LazyGraphClusterSeverityCounts v-if="type === 'severity'" :data="data" />
      <LazyGraphClusterResultCounts v-else :data="data" />
    </v-card>
  </app-row>
</template>

<script setup lang="ts">
import type { ClusterCharts, ViewType } from "~/types/core";
const props = defineProps<{ data: ClusterCharts; source: string; hideCluster: boolean; type: ViewType }>();
</script>
