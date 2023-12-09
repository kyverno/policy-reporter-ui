<template>
  <static-autocomplete
    label="Severity"
    query-name="severities"
    :items="[Severity.LOW, Severity.MEDIUM, Severity.HIGH]"
    v-bind="$attrs"
  />
</template>

<script lang="ts" setup>
import StaticAutocomplete from "~/components/StaticAutocomplete.vue";
import { Severity } from "../../types";
</script>
