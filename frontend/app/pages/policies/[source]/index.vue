<template>
  <page-layout v-if="sources"
               :title="capilize(source)"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="source"
  >
    <resource-scroller :list="sources">
      <template #default="{ item }">
        <policy-source-group :source="item" />
      </template>
    </resource-scroller>
  </page-layout>
  <unauthorized v-if="error?.status === 401" />
</template>

<script setup lang="ts">
const route = useRoute()
const source = computed(() => route.params.source as string)

const { load } = useSourceStore(source.value)
await load(source.value)

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ sources: [source.value], kinds: [...kinds.value, ...clusterKinds.value] }))

const { data: sources, refresh, error } = useAPI((api) => api.policySources(filter.value))

watch(filter, onChange(refresh))
</script>
