<template>
  <v-divider />
  <v-list-item :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.policy } }">
    <template v-slot:prepend>
      <v-btn v-if="details" class="mr-2" variant="text" :icon="!open ? `mdi-chevron-up` : `mdi-chevron-down`" @click.stop.prevent="open = !open"></v-btn>
      <AvatarSeverity :severity="item.severity ?? Severity.INFO" />
    </template>
    <v-list-item-title>
      {{ item.policy }}
    </v-list-item-title>
    <template v-slot:append>
      <ResultChip :status="Status.PASS" :count="item.results.pass" tooltip="pass results" />
      <ResultChip class="ml-2" :status="Status.WARN" :count="item.results.warn" tooltip="warning results" />
      <ResultChip class="ml-2" :status="Status.FAIL" :count="item.results.fail" tooltip="fail results" />
      <ResultChip class="ml-2" :status="Status.ERROR" :count="item.results.error" tooltip="error results" />
    </template>
  </v-list-item>
</template>

<script setup lang="ts">
import { type PropType } from "vue";
import { type PolicyResult, type Filter, Status, Severity } from "~/modules/core/types";

const open = ref(false)

const props = defineProps({
  item: { type: Object as PropType<PolicyResult>, required: true },
  details: { type: Boolean, default: false },
  filter: { type: Object as PropType<Filter>, required: false },
})
</script>
