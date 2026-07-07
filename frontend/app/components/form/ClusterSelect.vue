<template>
  <v-select
    class="mr-2"
    :model-value="cluster"
    :items="clusters"
    item-title="title"
    item-value="value"
    variant="outlined"
    hide-details
    density="compact"
    prepend-inner-icon="mdi-kubernetes"
    style="min-width: 140px;"
    @update:model-value="select"
    v-if="clusters.length > 1"
  />
</template>

<script lang="ts" setup>
import { useConfigStore } from "~/store/config";
import { cluster } from "~/core/api";

const router = useRouter()
const route = useRoute()

const store = useConfigStore()
const { $coreAPI } = useNuxtApp()

const clusters = computed(() => store.clusters.map(c => ({ title: c.name, value: c.slug })))

const select = (slug: string) => {
  $coreAPI.setPrefix(slug)
  router.push({ name: route.name as string, query: { ...route.query, cluster: slug }, params: route.params })
}

</script>
