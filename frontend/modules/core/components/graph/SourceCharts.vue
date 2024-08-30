<template>
  <GraphSourceStatus v-if="data.singleSource" :data="data" :source="data.sources[0]" />
  <template v-else>
    <GraphSourcesSeverities v-if="data.type === 'severity'" :data="data" :hide-cluster="hideCluster" />
    <GraphSourcesStatus v-else :data="data" :hide-cluster="hideCluster" />
  </template>
</template>

<script setup lang="ts">
import type {Dashboard} from "~/modules/core/types";

defineProps<{ data: Dashboard; hideCluster?: boolean }>();
</script>