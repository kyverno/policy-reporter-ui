<template>
  <app-autocomplete
    :source="source"
    label="Namespace"
    query-name="namespaces"
    cluster-api="namespaces"
    v-bind="$attrs"
  />
</template>

<script lang="ts" setup>
import AppAutocomplete from "~~/components/AppAutocomplete.vue";

defineProps({
  source: { type: String, default: undefined },
});
</script>
