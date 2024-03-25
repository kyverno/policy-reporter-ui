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
    <ResultChip v-if="showSkipped" class="ml-2" :status="Status.SKIP" :count="item[Status.SKIP]" tooltip="skip results" />
    <template v-slot:append>
      <ResultChip v-for="status in showed" :key="status" class="ml-2" :status="status" :count="item[status]" :tooltip="`${status} results`" />

      <exception-dialog v-if="source && exceptions" :resource="item.id" :source="source" :height="32" />
    </template>
  </v-list-item>
  <resource-source-results v-if="open" :id="item.id" :filter="filter" />
</template>

<script setup lang="ts">
import { type Filter, type ResourceResult, Status } from '~/modules/core/types'
import { type PropType } from "vue";
import { useStatusInjection } from "~/composables/status";
import ExceptionDialog from "~/modules/core/components/resource/ExceptionDialog.vue";
import {injectSourceContext} from "~/composables/source";

const open = ref(false)

const props = defineProps({
  item: { type: Object as PropType<ResourceResult>, required: true },
  details: { type: Boolean, default: false },
  exceptions: { type: Boolean, default: false },
  filter: { type: Object as PropType<Filter>, required: false },
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

const showSkipped = computed(() => status.value.includes(Status.SKIP) && !!props.item?.[Status.SKIP])
const showed = computed(() => status.value.filter((s) => s !== Status.SKIP))
</script>
