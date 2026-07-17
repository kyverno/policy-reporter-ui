<template>
  <page-layout v-if="data"
               title="Policy Dashboard"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
  >
    <policy-source-group v-for="item in data.sources" :key="item.name" :source="item" />
  </page-layout>
  <unauthorized v-if="error?.status === 401" />
</template>

<script setup lang="ts">
const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ kinds: [...kinds.value, ...clusterKinds.value] }))

const { data, refresh, error } = useAPI((api) => api.policySources(filter.value))

watch(filter, onChange(refresh))

useNamespacedKindProvider(data)
useClusterKindProvider(data)
</script>