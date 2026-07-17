<template>
  <page-layout v-if="data"
               :title="capilize(source)"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="source"
  >
    <resource-scroller :list="data.sources" class="mt-6">
      <template #default="{ item }">
        <policy-source-group :source="item" />
      </template>
    </resource-scroller>
  </page-layout>
  <unauthorized v-if="error?.status === 401" />
</template>

<script setup lang="ts">
import type { Filter } from '~/types/core'

const route = useRoute()

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const source = computed<string>(() => route.params.source as string)

const filter = computed((): Filter => ({
    sources: source.value ? [source.value] : undefined,
}))

const { data, refresh, error } = useAPI((api) => api.policySources(filter.value))

watch(filter, onChange(refresh))

useNamespacedKindProvider(data)
useClusterKindProvider(data)
</script>
