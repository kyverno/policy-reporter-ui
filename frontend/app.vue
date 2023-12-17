<template>
  <NuxtLayout :key="cluster">
    <NuxtPage :sources="sources" />
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { useAPI } from "~/modules/core/composables/api";
import { cluster } from "~/modules/core/api";

const { data: sources } = useAPI(
    (api) => api.sources().then((source) => source.map(s => s.name)),
    { default: () => [] }
);

</script>

<style>
  .main-height {
    min-height: calc(100vh - 64px);
  }

  .no-scrollbar {
    overflow: hidden!important;
  }
</style>
