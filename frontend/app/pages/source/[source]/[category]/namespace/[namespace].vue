<template>
  <page-layout
    :title="'Namespace Dashboard'"
    v-model:kinds="kinds"
    :source="source"
    v-if="namespace"
  >
    <template #append>
        <v-btn variant="text" class="mr-3" color="white" prepend-icon="mdi-arrow-left" :to="{ name: 'source-source-category', params: { source, category }, query: { ...route.query, mode: 'compact' } }">back</v-btn>
    </template>
    <resource-list :namespace="namespace" :details="data.multiSource" :exceptions="data.exceptions" :per-page="100" />
  </page-layout>
</template>

<script setup lang="ts">
import { APIFilter } from '~/provider/dashboard';

const route = useRoute();

const namespace = computed(() => route.params.namespace as string);

const { kinds, filter, source, category } = useFilter()

const { data, refresh } = useAPI(api => api.namespace(filter.value))

watch(filter, onChange(refresh))

const store = useSourceStore(source.value);
await store.load(source.value);

provide(APIFilter, filter);

useSourceContext(source);
useStatusProvider(data)
useSeveritiesProvider(data)
</script>
