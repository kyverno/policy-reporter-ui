<template>
  <v-select
    label="Status"
    :items="[Status.FAIL, Status.WARN, Status.ERROR, Status.PASS]"
    hide-details
    v-bind="$attrs"
  />
</template>

<script lang="ts" setup>
import { Status } from "../../types";
</script>
