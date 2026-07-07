<template>
  <app-row class="mb-2">
    <v-card>
      <v-toolbar color="header">
        <v-toolbar-title>Clusters</v-toolbar-title>
      </v-toolbar>
    </v-card>
  </app-row>
  <slot :clusters="list">
    <resource-scroller :list="list">
      <template #default="{ item }">
        <resource-total-list :cluster="item" :details="false" />
      </template>
    </resource-scroller>
  </slot>
</template>

<script setup lang="ts">
const props = defineProps<{ clusters: Array<{ name: string; slug: string }>; }>()

const list = computed(() => {
  return props.clusters || []
})

</script>
