<template>
  <app-autocomplete
    :source="source"
    :namespaced="namespaced"
    label="Policy"
    query-name="policies"
    cluster-api="clusterPolicies"
    namespaced-api="namespacedPolicies"
    v-bind="$attrs"
  />
</template>

<script lang="ts" setup>
import AppAutocomplete from "~/components/AppAutocomplete.vue";

const props = defineProps({
  namespaced: { type: Boolean, default: false },
  source: { type: String, default: undefined },
});
</script>
