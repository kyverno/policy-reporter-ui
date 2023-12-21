<template>
  <template v-for="item in data" :key="item.source">
    <v-divider />
    <v-list-item class="bg-grey-lighten-4">
      <v-list-item-title>{{ capilize(item.source) }}</v-list-item-title>
      <template v-slot:append>
        <ResultChip :status="Status.PASS" :count="item.pass" tooltip="pass results" />
        <ResultChip class="ml-2" :status="Status.WARN" :count="item.warn" tooltip="warning results" />
        <ResultChip class="ml-2" :status="Status.FAIL" :count="item.fail" tooltip="fail results" />
        <ResultChip class="ml-2" :status="Status.ERROR" :count="item.error" tooltip="error results" />
      </template>
    </v-list-item>
  </template>
</template>

<script setup lang="ts">
import { type Filter, Status } from '../types'
import { capilize } from "../layouthHelper";

const props = defineProps<{ id: string; filter?: Filter; }>()

const { data } = useAPI(
    (api) => api.resourceResults(props.id, props.filter),
    {
      default: () => [],
    }
);

</script>
