<template>
  <v-chip :color="`severity-${severity}`" dark v-bind="$attrs">
    {{ severity }}
  </v-chip>
</template>

<script lang="ts" setup>
import type { Severity } from "~/modules/core/types";

defineProps<{ severity: Severity }>()
</script>
