<template>
    <v-chip :prepend-icon="`mdi-${icons[status]}`" :color="`status-${status}`" variant="flat">
        <div class="text-right" style="min-width: 40px;">{{ count || 0 }}</div>
        <v-tooltip :text="tooltip" location="top" activator="parent" />
    </v-chip>
</template>

<script setup lang="ts">
import { Status } from '../types'

defineProps<{
  status: Status;
  count?: number;
  tooltip: string;
}>()

const icons = {
    [Status.PASS]: 'check',
    [Status.WARN]: 'alert',
    [Status.FAIL]: 'alert-circle',
    [Status.ERROR]: 'close',
    [Status.SKIP]: 'slash-forward',
    [Status.SUMMARY]: 'alert',
}
</script>
