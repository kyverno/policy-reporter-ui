<template v-if="data">
  <v-row>
    <v-col cols="12" md="4">
      <v-card>
        <v-card-text>
          <GraphStatusDistribution :data="data.charts.findings" :title="category" />
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12" md="8">
      <v-card style="height: 100%;">
        <v-card-text style="height: 100%" v-if="showExpanded">
          <GraphStatusPerNamespace :data="data.charts.namespaceScope[source].complete" />
        </v-card-text>
        <v-card-text style="height: 100%" v-else>
          <GraphStatusPerNamespace :data="data.charts.namespaceScope[source].preview" />
        </v-card-text>
        <v-btn v-if="hasPreview" variant="outlined" size="small" @click="expand = !expand" style="position: absolute; bottom: 10px; right: 10px;" rounded="0">
          <span v-if="showExpanded">Show preview</span>
          <span v-else>Show Complete List</span>
        </v-btn>
      </v-card>
    </v-col>
  </v-row>
  <template v-if="data.clusterScope">
    <app-row>
      <v-card :title="`${capilize(source)} Cluster Scoped Results`">
        <GraphClusterResultCounts :data="data.charts.clusterScope[source]" />
      </v-card>
    </app-row>
  </template>
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import type { Dashboard } from "~/modules/core/types";

const props = defineProps<{ data: Dashboard; source: string; category?: string; }>();

console.log(props.data)

const expand = ref(false)

const hasPreview = computed(() => !!props.data.charts.namespaceScope[props.source].preview)

const showExpanded = computed(() => {
  if (!props.data.charts.namespaceScope[props.source].preview) {
    return true
  }

  return expand.value
})
</script>
