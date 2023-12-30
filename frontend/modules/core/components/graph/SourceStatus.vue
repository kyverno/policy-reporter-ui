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
      <v-card style="height: 100%">
        <v-card-text style="height: 100%">
          <GraphStatusPerNamespace :data="data.charts.namespaceScope[source]" />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
  <template v-if="data.clusterScope">
    <v-row>
      <v-col>
        <v-card :title="`${capilize(source)} cluster scoped results`">
          <GraphClusterResultCounts :data="data.charts.clusterScope[source]" />
        </v-card>
      </v-col>
    </v-row>
  </template>
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import type { Dashboard } from "~/modules/core/types";

const props = defineProps<{ data: Dashboard<false>, source: string; category?: string; }>();
</script>
