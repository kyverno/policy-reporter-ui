<template>
  <v-avatar :color="`severity-${severity}`" dark v-bind="$attrs" variant="tonal">
    {{ severity.charAt(0).toUpperCase() }}
    <v-tooltip :text="`${severity} severity`" location="top" activator="parent" />
  </v-avatar>
</template>

<script lang="ts" setup>
import type { Severity } from "~/types/core";

const props = defineProps<{ severity: Severity }>()
</script>
