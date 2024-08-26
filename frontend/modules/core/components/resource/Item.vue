<template>
  <v-divider />
  <v-list-item :to="{ name: 'resource-id', query: rsFilter, params: { id: item.id }}">
    <template v-if="details" v-slot:prepend>
      <v-btn class="mr-2" variant="text" :icon="!open ? `mdi-chevron-up` : `mdi-chevron-down`" @click.stop.prevent="open = !open"></v-btn>
    </template>
    <v-list-item-title>
      {{ item.name }}
    </v-list-item-title>
    <v-list-item-subtitle>{{ item.apiVersion }} {{ item.kind }}</v-list-item-subtitle>
    <template v-slot:append>
      <ResultChip v-if="showSkipped" class="ml-2" :status="Status.SKIP" :count="item[Status.SKIP]" tooltip="skip results" />
      <ResultChip v-for="status in showed" :key="status" class="ml-2" :status="status" :count="item[status]" :tooltip="`${status} results`" />

      <resource-exception-dialog v-if="source && exceptions" :resource="item.id" :source="source" :category="category" :height="32" btn-class="ml-4" />
    </template>
  </v-list-item>
  <resource-source-results v-if="open" :id="item.id" :filter="filter" :show-skipped="showSkipped" :showed="showed" />
</template>

<script setup lang="ts">
import { type Filter, type ResourceResult, Status } from '~/modules/core/types'
import { type PropType } from "vue";
import { useStatusInjection } from "~/composables/status";
import {injectSourceContext} from "~/composables/source";

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
const source = injectSourceContext()

const showed = computed(() => status.value.filter((s) => s !== Status.SKIP))

const category = computed(() => {
  if (props?.filter?.categories?.length !== 1) return undefined

  return props.filter?.categories[0]
})
</script>
