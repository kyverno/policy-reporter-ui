<template>
  <page-layout
    title="Namespace Dashboard"
    hide-report
    v-model:kinds="kinds"
    :source="source"
    v-if="data"
  >
    <template #append>
        <v-btn variant="text" class="mr-3" color="white" prepend-icon="mdi-arrow-left" :to="{ name: 'custom-boards-id', params: { id }, query: { ...route.query, mode: 'compact' } }">back</v-btn>
    </template>
    <custom-board-table v-if="showResults" :namespace="namespace" :sources="data.sources" :id="id" />
    <custom-board-list v-else :namespace="namespace" :details="data.multiSource" :id="id" :exceptions="data.exceptions" :per-page="20" />
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
const { showResults } = useDashboardHelper(data)

const source = computed(() => data.value.singleSource ? data.value.sources[0] : undefined)

watch(filter, onChange(refresh))

provide(APIFilter, computed(() => ({
  ...filter.value,
  namespaces: [namespace.value],
})))

useDashboardProvider(data)
</script>
