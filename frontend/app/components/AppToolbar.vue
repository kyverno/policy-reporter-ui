<template>
<v-toolbar color="secondary">
    <template #title>
        {{ title }}
    </template>
    <template #append>
    <search :model-value="searchText" @update:model-value="(value) => emit('update:searchText', value)" class="mr-4" style="min-width: 400px; float: left;" />
    <v-btn :icon="open ? 'mdi-chevron-down' : 'mdi-chevron-up'" @click="() => emit('update:open', !open)" variant="text" />
    </template>
</v-toolbar>
</template>

<script lang="ts" setup>
const props = defineProps<{
    title: string;
    searchText: string;
    open: boolean;
}>()

const emit = defineEmits<{
    (e: 'update:searchText', value: string): void
    (e: 'update:open', value: boolean): void
}>()
</script>
