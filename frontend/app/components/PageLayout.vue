<template>
  <v-container fluid class="py-4 px-4 main-height">
    <app-row>
      <v-card>
        <v-toolbar color="header" elevation="2">
          <v-toolbar-title v-if="title">{{ title }}</v-toolbar-title>
          <template #append v-if="mdAndUp">
            <slot name="prepend" />
            <policy-report-dialog :source="store || source" :category="category" v-if="source && !hideReport" />
            <FormModeSelect style="min-width: 152px; max-width: 100%; margin-right: 15px;" v-model="mode" v-if="mode" />
            <FormKindAutocomplete style="min-width: 300px; max-width: 100%; margin-right: 15px;" v-model="kinds" :source="store || source" />
            <FormClusterKindAutocomplete v-if="!nsScoped" style="min-width: 300px; margin-right: 15px;" v-model="clusterKinds" :source="store || source" />
            <slot name="append" />
          </template>
        </v-toolbar>
        <v-container fluid v-if="!mdAndUp">
          <app-row v-if="source && !hideReport">
            <policy-report-dialog :source="store || source" :category="category" block />
          </app-row>
          <app-row v-if="mode">
            <FormModeSelect style="width: 100%;" v-model="mode" />
          </app-row>
          <app-row>
            <FormKindAutocomplete style="width: 100%;" v-model="kinds" :source="store || source" />
          </app-row>
          <app-row v-if="!nsScoped">
            <FormClusterKindAutocomplete style="width: 100%;" v-model="clusterKinds" :source="store || source" />
          </app-row>
        </v-container>
      </v-card>
    </app-row>
    <slot />
  </v-container>
</template>

<script setup lang="ts">
import { useDisplay } from "vuetify";
import { ClusterKinds, NamespacedKinds } from "~/provider/dashboard";
import type { Mode } from "~/types/core";
const { mdAndUp } = useDisplay()

const props = defineProps<{ title?: string; category?: string; source?: string; store?: string; nsScoped?: boolean; kinds?: string[]; clusterKinds?: string[]; hideReport?: boolean; mode?: Mode }>()

const kinds = ref<string[]>(props.kinds ?? [])
const clusterKinds = ref<string[]>(props.clusterKinds ?? [])
const mode = ref<Mode>(props.mode ?? '')

provide(NamespacedKinds, kinds)
provide(ClusterKinds, clusterKinds)

const emit = defineEmits<{
  'update:kinds': [kinds: string[]];
  'update:clusterKinds': [kinds: string[]];
  'update:mode': [mode: Mode];
}>()

watch(kinds, (k) => emit('update:kinds', [...k]))
watch(clusterKinds, (k) => emit('update:clusterKinds', [...k]))
watch(mode, (m) => emit('update:mode', m), { immediate: true })
</script>
