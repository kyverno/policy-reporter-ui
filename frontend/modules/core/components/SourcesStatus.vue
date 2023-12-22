<template>
  <v-row>
    <v-col cols="12" md="6">
      <v-card>
        <v-card-text class="mb-0 pb-0">
          <ChartFindings :findings="data" :status="status1" :key="status1" />
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
            <v-btn
                :color="mapStatus(Status.PASS)"
                clss="py-1"
                size="small"
                v-if="totals[Status.PASS]"
                :value="Status.PASS"
            >{{ Status.PASS }}</v-btn
            >
            <v-btn
                :color="mapStatus(Status.SKIP)"
                clss="py-1"
                size="small"
                v-if="totals[Status.SKIP]"
                :value="Status.SKIP"
            >{{ Status.SKIP }}</v-btn
            >
          </v-btn-toggle>
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col cols="12" md="6">
      <v-card>
        <v-card-text class="mb-0 pb-0">
          <ChartFindings :findings="data" :status="status2" :key="status2" />
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
            <v-btn
                :color="mapStatus(Status.FAIL)"
                clss="py-1"
                size="small"
                v-if="totals[Status.FAIL]"
                :value="Status.FAIL"
            >{{ Status.FAIL }}</v-btn
            >
            <v-btn
                :color="mapStatus(Status.WARN)"
                clss="py-1"
                size="small"
                v-if="totals[Status.WARN]"
                :value="Status.WARN"
            >{{ Status.WARN }}</v-btn
            >
            <v-btn
                :color="mapStatus(Status.ERROR)"
                clss="py-1"
                size="small"
                v-if="totals[Status.ERROR]"
                :value="Status.ERROR"
            >{{ Status.ERROR }}</v-btn
            >
          </v-btn-toggle>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
  <v-row v-if="sources">
    <v-col>
      <v-card title="Result Overview">
        <v-tabs v-model="source" bg-color="transparent" fixed-tabs>
          <v-tab v-for="item in items" :key="item.value" :value="item.value" rounded="0">
            {{ item.title }}
          </v-tab>
        </v-tabs>
        <v-divider />
        <template v-if="source">
          <wait :time="1000" :key="source">
            <v-card-text>
              <ChartStatusPerNamespace :source="source" />
            </v-card-text>
            <template v-if="!hideCluster">
              <v-divider />
              <v-card-title>
                Cluster Scoped Results
              </v-card-title>
              <v-card-text>
                <ChartClusterResultCounts :source="source" class="px-0 pb-0" />
              </v-card-text>
            </template>
            <template #placeholder>
              <v-card-text>
                <v-progress-linear indeterminate color="primary" />
              </v-card-text>
            </template>
          </wait>
        </template>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { mapStatus } from '../mapper';
import { type Filter, type FindingCounts, Status } from '../types';
import { capilize } from "~/modules/core/layouthHelper";
import { clusterKinds, kinds } from "~/modules/core/store/filter";

const props = defineProps<{ data: FindingCounts; filter?: Filter; hideCluster?: boolean }>();

const status1 = ref(Status.PASS);
const status2 = ref(Status.FAIL);

const initTotals = () => ({
  [Status.PASS]: 0,
  [Status.SKIP]: 0,
  [Status.FAIL]: 0,
  [Status.WARN]: 0,
  [Status.ERROR]: 0,
});

const totals = ref(initTotals());
const sources = ref<string[]>([]);

watch(() => props.data, (findings: FindingCounts) => {
  if (!findings) return;

  const results = initTotals();

  findings.counts.forEach((f) => {
    Object.keys(f.counts).forEach((s) => {
      results[s as Status] += f.counts[s as Status];
    });
  });

  totals.value = results;
  sources.value = [...new Set(findings.counts.map(c => c.source).sort((a, b) => a.localeCompare(b)))]
}, { immediate: true });

const source = ref('');

const items = computed(() => {
  if (!sources.value) return []

  return sources.value.map(s => ({
    title: capilize(s),
    value: s
  }))
})

watch(sources, (s) => {
  if (!s || !s.length || source.value) return

  source.value = s[0]
})
</script>
