<template>
  <v-divider />
  <v-list-item :to="{ name: 'resource-id-type', query: rsFilter, params: { id: item.id, type: viewType }}">
    <template v-if="details" v-slot:prepend>
      <v-btn class="mr-2" variant="text" :icon="!open ? `mdi-chevron-up` : `mdi-chevron-down`" @click.stop.prevent="open = !open"></v-btn>
    </template>
    <v-list-item-title>
      {{ item.name }}
    </v-list-item-title>
    <v-list-item-subtitle>{{ item.apiVersion }} {{ item.kind }}</v-list-item-subtitle>
    <template v-if="viewType === 'severity'" v-slot:append>
      <resource-severities-chips :items="item.severities as any" :id="item.id" :category="category" :exceptions="exceptions" />
    </template>
    <template v-else v-slot:append>
      <resource-status-chips :items="item.status as any" :id="item.id" :category="category" :exceptions="exceptions" :show-skipped="showSkipped" :showed="showed" />
    </template>
  </v-list-item>
  <resource-source-results v-if="open" :id="item.id" :filter="filter" :show-skipped="showSkipped" :showed="showed" />
</template>

<script setup lang="ts">
import { type Filter, type ResourceResult, Status } from '~/modules/core/types'
import { type PropType } from "vue";
import { useStatusInjection } from "~/composables/status";
import {injectDashboardType} from "~/composables/dashboard";

const open = ref(false)

const props = defineProps({
  item: { type: Object as PropType<ResourceResult>, required: true },
  details: { type: Boolean, default: false },
  exceptions: { type: Boolean, default: false },
  filter: { type: Object as PropType<Filter>, required: false },
  showSkipped: { type: Boolean, default: false },
})

const rsFilter = computed(() => {
  if (!props.filter) return undefined

  return {
    sources: props.filter.sources,
    categories: props.filter.categories,
  }
})

const status = useStatusInjection()
const viewType = injectDashboardType()

const showed = computed(() => status.value.filter((s) => s !== Status.SKIP))

const category = computed(() => {
  if (props?.filter?.categories?.length !== 1) return undefined

  return props.filter?.categories[0]
})
</script>
