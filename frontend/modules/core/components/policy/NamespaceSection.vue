<template>
  <v-row>
    <v-col>
      <v-card>
        <v-toolbar color="header">
          <v-toolbar-title>Namespace Scoped Results</v-toolbar-title>
          <template #append>
            <form-namespace-autocomplete v-model="internal" :items="props.namespaces" />
          </template>
        </v-toolbar>
      </v-card>
    </v-col>
  </v-row>
  <slot :namespaces="list">
    <resource-scroller :list="list">
      <template #default="{ item }">
        <app-row>
          <policy-results :namespace="item" :source="source" :policy="policy" />
        </app-row>
      </template>
    </resource-scroller>
  </slot>
</template>

<script setup lang="ts">
const props = defineProps<{ namespaces: string[]; source: string; policy?: string;  }>()

const internal = ref<string[]>([])

const list = computed(() => {
  if (internal.value.length > 0) return internal.value

  return props.namespaces || []
})

</script>
