<template v-if="data">
  <v-row>
    <v-col cols="12" md="4">
      <v-card>
        <v-card-text>
          <ChartStatusDistribution :findings="data" :title="policy" />
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12" md="8" v-if="source">
      <v-card style="height: 100%">
        <v-card-text style="height: 100%">
          <ChartStatusPerNamespace :title="policy" :source="source" :key="source" />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
  <template v-if="!hideCluster">
    <v-row>
      <v-col>
        <v-card :title="`${policy} cluster scoped results`">
          <ChartClusterResultCounts :source="source" />
        </v-card>
      </v-col>
    </v-row>
  </template>
</template>

<script setup lang="ts">
import type { SourceFindings } from "~/modules/core/types";

const props = defineProps<{ data: SourceFindings, hideCluster?: boolean; policy: string; }>();

const source = computed(() => props.data.source)
</script>
