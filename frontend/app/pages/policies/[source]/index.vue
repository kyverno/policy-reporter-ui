<template>
  <page-layout v-if="sources"
               :title="capilize(source)"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="source"
  >
    <resource-scroller :list="sources" class="mt-6">
      <template #default="{ item }">
        <policy-source-group :source="item" />
      </template>
    </resource-scroller>
  </page-layout>
  <unauthorized v-if="error?.status === 401" />
</template>

<script setup lang="ts">
const { filter, clusterKinds, kinds, source } = useFilter()

const { load } = useSourceStore(source.value)
await load(source.value)

const { data: sources, refresh, error } = useAPI((api) => api.policySources(filter.value))

watch(filter, onChange(refresh))
</script>
