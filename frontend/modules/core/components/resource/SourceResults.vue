<template>
  <template v-for="item in data" :key="item.source">
    <v-divider />
    <v-list-item :class="bg">
      <v-list-item-title>{{ capilize(item.source) }}</v-list-item-title>
      <template v-slot:append>
        <ResultChip v-if="showSkipped" class="ml-2" :status="Status.SKIP" :count="item[Status.SKIP]" tooltip="skip results" />
        <ResultChip v-for="status in showed" :key="status" class="ml-2" :status="status" :count="item[status]" :tooltip="`${status} results`" />
      </template>
    </v-list-item>
  </template>
</template>

<script setup lang="ts">
import { type Filter, Status } from '~/modules/core/types'
import { capilize } from "~/modules/core/layouthHelper";

const props = defineProps<{ id: string; filter?: Filter; showSkipped: boolean; showed: Status[]; }>()

const bg = useBGColor()

const { data } = useAPI(
    (api) => api.resourceResults(props.id, props.filter),
    {
      default: () => [],
    }
);

</script>
