<template>
  <app-autocomplete
    :source="source"
    :namespaced="namespaced"
    label="Rule"
    query-name="rules"
    cluster-api="clusterRules"
    namespaced-api="namespacedRules"
    v-bind="$attrs"
  />
</template>

<script lang="ts" setup>
import AppAutocomplete from "~/components/AppAutocomplete.vue";

defineProps({
  namespaced: { type: Boolean, default: false },
  source: { type: String, default: undefined },
});
</script>
