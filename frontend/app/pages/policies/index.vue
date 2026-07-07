<template>
  <page-layout v-if="sources"
               title="Policy Dashboard"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
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


const { load } = useSourceStore()
await load()

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ kinds: [...kinds.value, ...clusterKinds.value] }))

const { data: sources, refresh, error } = useAPI((api) => api.policySources(filter.value))

watch(filter, onChange(refresh))
</script>
