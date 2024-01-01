<template>
  <v-select
    :model-value="cluster"
    :items="clusters"
    item-title="title"
    item-value="value"
    variant="outlined"
    hide-details
    density="compact"
    prepend-inner-icon="mdi-kubernetes"
    style="max-width: 140px;"
    @update:model-value="input"
    v-if="clusters.length > 1"
  />
</template>

<script lang="ts" setup>
import { useConfigStore } from "~/store/config";
import type { CoreAPI } from "../../api";
import { cluster } from "../../api";

const store = useConfigStore()
const { $coreAPI } = useNuxtApp()

const clusters = computed(() => store.clusters.map(c => ({ title: c.name, value: c.slug })))

const input = (slug: string) => { ($coreAPI as CoreAPI).setPrefix(slug) }

</script>
