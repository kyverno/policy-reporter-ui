<template>
  <v-chip :color="`severity-${severity}`" dark v-bind="$attrs">
    {{ severity }}
  </v-chip>
</template>

<script lang="ts" setup>
import type { Severity } from "~/types/core";

defineProps<{ severity: Severity }>()
</script>
