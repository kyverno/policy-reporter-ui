<template>
  <NuxtLayout :key="cluster">
    <NuxtPage />
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { cluster } from "~/core/api";
const { $coreAPI } = useNuxtApp()

const router = useRouter()
const route = useRoute()

const start = computed(() => route.query.cluster as string || '')

if (start.value) {
  $coreAPI.setPrefix(start.value)
  router.push({ name: route.name as string, query: { ...route.query, cluster: undefined }, params: route.params })
}
</script>

<style>
  .main-height {
    min-height: calc(100vh - 64px);
  }

  .no-scrollbar {
    overflow: hidden!important;
  }

  .v-theme--dark .v-data-table-footer,
  .v-theme--dark .top-border {
    border-top: 1px solid rgba(255, 255, 255, 0.12)
  }

  .v-theme--light .v-data-table-footer,
  .v-theme--light .top-border {
    border-top: 1px solid #E1DCDF;
  }
</style>
