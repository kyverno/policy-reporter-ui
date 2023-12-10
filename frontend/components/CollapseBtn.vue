<template>
  <v-btn class="ml-2" variant="text" :icon="!modelValue ? `mdi-chevron-up` : `mdi-chevron-down`" @click.stop.prevent="$emit('update:modelValue', !modelValue)"></v-btn>
</template>

<script setup lang="ts">
defineProps<{ modelValue: boolean; }>()
</script>
