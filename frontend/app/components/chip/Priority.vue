<template>
  <v-chip :color="color" v-bind="$attrs">
    <template v-if="prefix">{{ prefix }}&nbsp;</template>{{ priority ? priority : Priority.DEBUG }}
  </v-chip>
</template>

<script lang="ts" setup>
import { Priority } from "~/types/core";
import { mapPriority } from "~/core/mapper";

const props = defineProps<{ priority?: Priority; prefix?: string; }>()

const color = mapPriority(props.priority)
</script>
