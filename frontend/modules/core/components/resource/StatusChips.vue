<template>
  <ResultChip v-if="showSkipped" class="ml-2" :status="Status.SKIP" :count="items[Status.SKIP]" tooltip="skip results" />
  <ResultChip v-for="status in showed" :key="status" class="ml-2" :status="status" :count="items[status]" :tooltip="`${status} results`" />
  <resource-exception-dialog v-if="source && exceptions" :resource="id" :source="source" :category="category" :height="32" btn-class="ml-4" />
</template>

<script setup lang="ts">
import { Status } from '~/modules/core/types'
import { type PropType } from "vue";
import {injectSourceContext} from "~/composables/source";

defineProps({
  id: { type: String, required: true },
  items: { type: Object as PropType<{ [status in Status]: number }>, required: true },
  exceptions: { type: Boolean, default: false },
  showSkipped: { type: Boolean, default: false },
  showed: { type: Array as PropType<Status[]>, default: false },
  category: { type: String, default: undefined },
})

const source = injectSourceContext()
</script>
