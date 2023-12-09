<template>
  <app-autocomplete
    :source="source"
    label="Category"
    query-name="categories"
    cluster-api="categories"
    v-bind="$attrs"
  />
</template>

<script lang="ts" setup>
import AppAutocomplete from "~/components/AppAutocomplete.vue";

const props = defineProps({
  source: { type: String, default: undefined },
});
</script>
