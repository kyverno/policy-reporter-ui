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
    <policy-scroller :list="list">
      <template #default="{ item }">
        <policy-results :namespace="item" :source="source" :policy="policy" :status="status" :exceptions="exceptions" />
      </template>
    </policy-scroller>
  </slot>
</template>

<script setup lang="ts">
import type { Status } from "~/modules/core/types";

const props = defineProps<{ namespaces: string[]; source: string; policy?: string; exceptions?: boolean; status?: Status[];  }>()

const internal = ref<string[]>([])

const list = computed(() => {
  if (internal.value.length > 0) return internal.value

  return props.namespaces || []
})

</script>
