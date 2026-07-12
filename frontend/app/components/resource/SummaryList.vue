<template>
  <v-card class="mb-4">
    <v-toolbar color="transparent">
      <v-toolbar-title>Namespace Summary</v-toolbar-title>
      <template #append>
        <Search v-model="search" style="min-width: 400px;" />
        <CollapseBtn v-model="open" />
      </template>
    </v-toolbar>
      <v-list v-if="data?.items?.length && open" lines="two">
        <resource-item v-for="item in filtered" :key="item.id" :item="item" :details="false" />
      </v-list>
      <template v-if="!data?.items?.length">
        <v-divider />
        <v-card-text>
          No resources for the selected kinds found
        </v-card-text>
      </template>
  </v-card>
</template>

<script setup lang="ts">
import { type ResourceResultList } from '~/types/core'
import CollapseBtn from "~/components/CollapseBtn.vue";

const props = defineProps<{ data: ResourceResultList; }>()

const search = ref('')
const open = ref(true)

const filtered = computed(() => {
  if (!search.value) {
    return props.data.items
  }
  return props.data.items.filter((item) => item.name.toLowerCase().includes(search.value.toLowerCase()))
})
</script>
