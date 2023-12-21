<template>
  <NuxtLayout :key="cluster">
    <NuxtPage />
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { cluster } from "~/modules/core/api";
</script>

<style>
  .main-height {
    min-height: calc(100vh - 64px);
  }

  .no-scrollbar {
    overflow: hidden!important;
  }
</style>
