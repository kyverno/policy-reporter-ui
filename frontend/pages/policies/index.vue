<template>
  <page-layout v-if="sources"
               title="Policy Dashboard"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
  >
    <resource-scroller :list="sources">
      <template #default="{ item }">
        <policy-source-group :source="item" />
      </template>
    </resource-scroller>
  </page-layout>
</template>

<script setup lang="ts">
import { onChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ kinds: [...kinds.value, ...clusterKinds.value] }))

const { data: sources, refresh } = useAPI((api) => api.policySources(filter.value))

watch(filter, onChange(refresh))
</script>
