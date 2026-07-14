<template>
  <page-layout
    :title="'Namespace Dashboard'"
    v-model:kinds="kinds"
    :source="source"
    v-if="data"
  >
    <template #append>
        <v-btn variant="text" class="mr-3" color="white" prepend-icon="mdi-arrow-left" :to="{ name: 'custom-boards-id', params: { id }, query: { ...route.query, mode: 'compact' } }">back</v-btn>
    </template>
    <resource-list :namespace="namespace" :details="data.multiSource" :exceptions="data.exceptions" :per-page="100" />
  </page-layout>
  <unauthorized v-if="error?.status === 401" />
</template>

<script setup lang="ts">
import { APIFilter } from '~/provider/dashboard';

const route = useRoute();

const namespace = computed(() => route.params.namespace as string);
const id = computed(() => route.params.id as string)

const { kinds, filter } = useFilter()

const { data, refresh, error } = useAPI((api) => api.customBoard(id.value, filter.value))

const source = computed(() => data.value.singleSource ? data.value.sources[0] : undefined)

const store = useSourceStore(id.value)

watchEffect(() => {
  if (!data.value) return;

  store.load(data.value.sources)

  if (data.value.namespaceKinds.length > 0) {
    kinds.value = data.value.namespaceKinds
  }
})

watch(filter, onChange(refresh))

provide(APIFilter, computed(() => ({
  ...filter.value,
  sources: data.value?.filterSources,
  namespaces: [namespace.value],
})))

useStatusProvider(data)
useSeveritiesProvider(data)
useDashboardType(data)
</script>
