<template>
    <v-chip :prepend-icon="`mdi-${icons[severity]}`" :color="`severity-${severity}`" variant="flat">
        <div class="text-right" style="min-width: 40px;">{{ count || 0 }}</div>
        <v-tooltip :text="tooltip" location="top" activator="parent" />
    </v-chip>
</template>

<script setup lang="ts">
import { Severity } from '../types'

defineProps<{
  severity: Severity;
  count?: number;
  tooltip: string;
}>()

const icons = {
    [Severity.UNKNOWN]: 'help',
    [Severity.INFO]: 'information',
    [Severity.LOW]: 'alert',
    [Severity.MEDIUM]: 'alert',
    [Severity.HIGH]: 'alert',
    [Severity.CRITICAL]: 'alert',
}
</script>
