<template>
  <SeverityChip v-for="severity in showed" :key="severity" class="ml-2" :severity="severity" :count="items[severity]" :tooltip="`${severity} results`" />
  <resource-exception-dialog v-if="source && exceptions" :resource="id" :source="source" :category="category" :height="32" btn-class="ml-4" />
</template>

<script setup lang="ts">
import {Severity} from '~/modules/core/types'
import { type PropType } from "vue";
import { injectSourceContext } from "~/composables/source";

defineProps({
  id: { type: String, required: true },
  items: { type: Object as PropType<{ [key in Severity]: number }>, required: true },
  exceptions: { type: Boolean, default: false },
  category: { type: String, required: false },
})

const showed = useSeveritiesInjection()
const source = injectSourceContext()
</script>
