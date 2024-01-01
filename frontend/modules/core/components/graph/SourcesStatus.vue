<template>
  <v-row>
    <v-col cols="12" md="6">
      <v-card>
        <v-card-text class="mb-0 pb-0">
          <GraphFindings :data="data.charts.findings[status1]" :status="status1" :key="status1" />
        </v-card-text>
        <v-card-actions style="margin-top: -40px">
          <v-btn-toggle
              v-model="status1"
              rounded="0"
              mandatory
              style="height: 40px"
              divided
              variant="outlined"
          >
            <status-btn v-if="data.total.perResult[Status.PASS]" :status="Status.PASS" />
            <status-btn v-if="data.total.perResult[Status.SKIP]" :status="Status.SKIP" />
          </v-btn-toggle>
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col cols="12" md="6">
      <v-card>
        <v-card-text class="mb-0 pb-0">
          <GraphFindings :data="data.charts.findings[status2]" :status="status2" :key="status2" :time="1200" />
        </v-card-text>
        <v-card-actions style="margin-top: -40px">
          <v-btn-toggle
              v-model="status2"
              rounded="0"
              mandatory
              style="height: 40px"
              divided
              variant="outlined"
          >
            <status-btn v-if="data.total.perResult[Status.FAIL]" :status="Status.FAIL" />
            <status-btn v-if="data.total.perResult[Status.WARN]" :status="Status.WARN" />
            <status-btn v-if="data.total.perResult[Status.ERROR]" :status="Status.ERROR" />
          </v-btn-toggle>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
  <app-row>
    <v-card title="Result Overview">
      <v-tabs v-model="source" bg-color="transparent" fixed-tabs>
        <v-tab v-for="item in data.sourcesNavi" :key="item.name" :value="item.name" rounded="0">
          {{ item.title }}
        </v-tab>
      </v-tabs>
      <v-divider />
      <template v-if="source">
        <v-card-text>
          <GraphStatusPerNamespace :data="data.charts.namespaceScope[source]" />
        </v-card-text>
        <template v-if="!hideCluster">
          <v-divider />
          <v-card-title>
            Cluster Scoped Results
          </v-card-title>
          <v-card-text>
            <GraphClusterResultCounts :data="data.charts.clusterScope[source]" class="px-0 pb-0" />
          </v-card-text>
        </template>
      </template>
    </v-card>
  </app-row>
</template>

<script setup lang="ts">
import StatusBtn from "~/components/StatusBtn.vue";
import { type Dashboard, Status } from "~/modules/core/types";

const props = defineProps<{ data: Dashboard<true>; hideCluster?: boolean }>();

const status1 = ref(Status.PASS);
const status2 = ref(Status.FAIL);

const source = ref('');

watch(() => props.data.sources, (s) => {
  if (!s || !s.length || source.value) return

  source.value = s[0]
})
</script>
