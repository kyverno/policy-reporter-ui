<template v-if="data">
  <v-row>
    <v-col cols="12" md="4">
      <v-card>
        <v-card-text>
          <GraphStatusDistribution :data="data.charts.findings" />
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12" md="8">
      <v-card style="height: 100%">
        <v-card-text style="height: 100%; position: relative;">
          <GraphStatusPerNamespace v-if="showExpanded" :data="data.charts.namespaceScope.complete" />
          <GraphStatusPerNamespace v-else :data="data.charts.namespaceScope.preview" />
          <v-btn v-if="hasPreview" variant="outlined" size="small" @click="expand = !expand" style="position: absolute; bottom: 10px; right: 10px;" rounded="0">
            <span v-if="showExpanded">Show preview</span>
            <span v-else>Show Complete List</span>
          </v-btn>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
  <template v-if="clusterScope">
    <app-row>
      <v-card :title="`${data.title} cluster scoped results`">
        <GraphClusterResultCounts :data="data.charts.clusterScope" />
      </v-card>
    </app-row>
  </template>
</template>

<script setup lang="ts">
import type { PolicyDetails } from "~/modules/core/types";

const props = defineProps<{ data: PolicyDetails, hideCluster?: boolean; policy: string; }>();

const expand = ref(false)

const hasPreview = computed(() => !!props.data.charts.namespaceScope.preview)

const showExpanded = computed(() => {
  if (!props.data.charts.namespaceScope.preview) {
    return true
  }

  return expand.value
})

const clusterScope = computed(() => {
  if (props.hideCluster) return false

  return Object.values(props.data.charts.clusterScope || {}).reduce((acc, res) => acc + res, 0) > 0
})
</script>
