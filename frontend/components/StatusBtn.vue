<template>
  <v-btn :color="`btn-${status}`" clss="py-1" size="small" :value="status">{{ status }}</v-btn>
</template>

<script setup lang="ts">
import { Status } from "~/modules/core/types";

const props = defineProps<{ status: Status }>()
</script>
