<template>
  <v-container fluid class="py-4 px-4 main-height">
    <app-row>
      <v-card>
        <v-toolbar color="header" elevation="2">
          <v-toolbar-title v-if="title">{{ title }}</v-toolbar-title>
          <template #append>
            <FormKindAutocomplete style="min-width: 300px; max-width: 100%; margin-right: 15px;" v-model="kinds" :source="source" />
            <FormClusterKindAutocomplete v-if="!nsScoped" style="min-width: 300px;" v-model="clusterKinds" :source="source" />
          </template>
        </v-toolbar>
      </v-card>
    </app-row>
    <slot />
  </v-container>
</template>

<script setup lang="ts">
import { ClusterKinds, NamespacedKinds } from "~/modules/core/provider/dashboard";

const props = defineProps<{ title?: string; source?: string; nsScoped?: boolean; kinds?: string[]; clusterKinds?: string[] }>()

const kinds = ref<string[]>(props.kinds ?? [])
const clusterKinds = ref<string[]>(props.clusterKinds ?? [])

provide(NamespacedKinds, kinds)
provide(ClusterKinds, clusterKinds)

const emit = defineEmits<{
  'update:kinds': [kinds: string[]];
  'update:clusterKinds': [kinds: string[]];
}>()

watch(kinds, (k) => emit('update:kinds', [...k]))
watch(clusterKinds, (k) => emit('update:clusterKinds', [...k]))
</script>
