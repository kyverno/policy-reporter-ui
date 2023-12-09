<template v-if="data">
  <v-row>
    <v-col cols="12" md="4">
      <v-card>
        <v-card-text>
          <ChartStatusDistribution :findings="data" />
        </v-card-text>
      </v-card>
    </v-col>
    <v-col cols="12" md="8" v-if="source">
      <v-card style="height: 100%">
        <v-card-text style="height: 100%">
          <ChartStatusPerNamespace :source="source" :key="source" />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
  <template v-if="statusCounts">
    <v-row>
      <v-col>
        <v-card :title="`${capilize(source)} cluster scoped results`">
          <ChartClusterResultCounts />
        </v-card>
      </v-col>
    </v-row>
  </template>
</template>

<script setup lang="ts">
import type { SourceFindings } from "~/modules/core/types";
import type { PropType } from "vue";
import { Status } from "~/modules/core/types";
import { mapStatus } from "~/modules/core/mapper";
import { capilize } from "~/modules/core/layouthHelper";

const props = defineProps<{ data: SourceFindings }>();

const source = computed(() => props.data.source)

const { data: sc } = useAPI(
    (api) => api.statusCount(), {
      default: () => [
        { status: Status.PASS, count: 0 },
        { status: Status.WARN, count: 0 },
        { status: Status.FAIL, count: 0 },
        { status: Status.ERROR, count: 0 },
      ],
    }
);

const statusCounts = computed<{ [status in Status]: number }>(() => {
  return sc.value.reduce((acc, item) => {
    if (item.status === Status.SKIP) return acc;

    acc[item.status] = item.count

    return acc
  }, {
    [Status.PASS]: 0,
    [Status.WARN]: 0,
    [Status.FAIL]: 0,
    [Status.ERROR]: 0,
  })
})
</script>
