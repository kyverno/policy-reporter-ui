<template>
  <v-chip :color="`severity-${severity}`" dark v-bind="$attrs">
    {{ severity }}
  </v-chip>
</template>

<script lang="ts" setup>
import type { Severity } from "~/modules/core/types";

const props = defineProps<{ severity: Severity }>()
</script>
