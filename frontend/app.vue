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

  .v-theme--dark .v-data-table-footer {
    border-top: 1px solid rgba(255, 255, 255, 0.12)
  }

  .v-theme--light .v-data-table-footer {
    border-top: 1px solid #E1DCDF;
  }
</style>
