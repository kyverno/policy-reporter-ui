<template>
  <v-chip :color="color" dark v-bind="$attrs">
    {{ severity }}
  </v-chip>
</template>

<script lang="ts" setup>
import type { PropType } from "vue";
import { Severity, Status } from "~/modules/core/types";

const props = defineProps({
  severity: { type: String as PropType<string>, default: null },
})

const color = computed(() => {
  switch (props.severity) {
    case Severity.INFO:
      return 'info lighten-1'
    case Severity.LOW:
      return 'info'
    case Severity.MEDIUM:
      return 'warning'
    case Severity.HIGH:
      return 'error'
    case Severity.CRITICAL:
      return 'error darken-2'
    default:
      return 'grey'
  }
})
</script>
