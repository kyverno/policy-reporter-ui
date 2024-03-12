<template>
  <v-divider />
  <v-list-item :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }}">
    <template v-slot:prepend>
      <CollapseBtn v-if="item.description" btn-class="mr-2" v-model="open" :size="40" />
      <AvatarSeverity :severity="item.severity ?? Severity.INFO" />
    </template>
    <v-list-item-title>
      {{ item.title }}
    </v-list-item-title>
    <template v-slot:append>
      <ResultChip :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }, query: { status: Status.PASS }}" :status="Status.PASS" :count="item.results.pass" tooltip="pass results" />
      <ResultChip :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }, query: { status: Status.WARN }}" class="ml-2" :status="Status.WARN" :count="item.results.warn" tooltip="warning results" />
      <ResultChip :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }, query: { status: Status.FAIL }}" class="ml-2" :status="Status.FAIL" :count="item.results.fail" tooltip="fail results" />
      <ResultChip :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }, query: { status: Status.ERROR }}" class="ml-2" :status="Status.ERROR" :count="item.results.error" tooltip="error results" />
    </template>
  </v-list-item>
  <template v-if="open">
    <v-divider />
    <v-list-item :class="`${bg} text-pre-line`">
      {{ item.description }}
    </v-list-item>
  </template>
</template>

<script setup lang="ts">
import { type PropType } from "vue";
import { type PolicyResult, type Filter, Status, Severity } from "~/modules/core/types";
import { capilize } from "~/modules/core/layouthHelper";

const open = ref(false)

const bg = useBGColor()

const props = defineProps({
  item: { type: Object as PropType<PolicyResult>, required: true },
  details: { type: Boolean, default: false },
  filter: { type: Object as PropType<Filter>, required: false },
})
</script>
