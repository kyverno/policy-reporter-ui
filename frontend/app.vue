<template>
  <NuxtLayout :key="cluster">
    <NuxtPage :sources="sources" />
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { useAPI } from "~/modules/core/composables/api";
import { cluster } from "~/modules/core/api";

const { data: sources } = useAPI<string[]>(
    (api) => Promise.all([
      api.namespacedSources(),
      api.clusterSources(),
    ]).then(([ns, cluster]) => [...new Set([...ns, ...cluster])]),
    {
      default: () => [],
    }
);

</script>

<style>
  .main-height {
    min-height: calc(100vh - 64px);
  }
</style>
