<template>
  <v-btn :class="btnClass ? btnClass : 'ml-2'" color="primary" variant="text" :icon="!modelValue ? `mdi-chevron-up` : `mdi-chevron-down`" @click.stop.prevent="$emit('update:modelValue', !modelValue)"></v-btn>
</template>

<script setup lang="ts">
defineProps<{ modelValue: boolean; btnClass?: string; }>()
</script>
