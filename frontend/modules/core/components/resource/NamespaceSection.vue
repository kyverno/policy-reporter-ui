<template>
  <app-row>
    <v-card>
      <v-toolbar color="header">
        <v-toolbar-title>Namespace Scoped Resources</v-toolbar-title>
        <template #append>
          <form-namespace-autocomplete v-model="internal" :items="props.namespaces" />
        </template>
      </v-toolbar>
    </v-card>
  </app-row>
  <slot :namespaces="list">
    <resource-scroller :list="list">
      <template #default="{ item }">
        <resource-list :namespace="item" :details="false" />
      </template>
    </resource-scroller>
  </slot>
</template>

<script setup lang="ts">
const props = defineProps<{ namespaces: string[]; }>()

const internal = ref<string[]>([])

const list = computed(() => {
  if (internal.value.length > 0) return internal.value

  return props.namespaces || []
})

</script>
