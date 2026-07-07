<template>
  <v-btn :color="`btn-${severity}`" clss="py-1" size="small" :value="severity">{{ severity }}</v-btn>
</template>

<script setup lang="ts">
import { Severity } from "~/types/core";

defineProps<{ severity: Severity }>()
</script>
