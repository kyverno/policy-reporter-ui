<template>
  <v-chip :color="`btn-${status}`" dark v-bind="$attrs" variant="tonal">
    {{ status }}
  </v-chip>
</template>

<script lang="ts" setup>
import { Status } from "~/types/core";

const props = defineProps<{ status: Status }>()
</script>
