<template>
  <v-row>
    <v-col cols="12" md="6">
      <v-card>
        <v-card-text class="mb-0 pb-0">
          <GraphFindings :data="(data.charts.findings[status1] as Findings)" :status="status1" :key="status1" />
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
          <GraphFindings :data="(data.charts.findings[status2] as Findings)" :status="status2" :key="status2" :time="1200" />
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
      <template  v-if="data.sourcesNavi.length === 0">
        <v-divider />
        <v-card-text>
          <v-alert type="error" variant="outlined">No results found</v-alert>
        </v-card-text>
      </template>
      <v-tabs v-model="source" bg-color="transparent" fixed-tabs v-else>
        <v-tab v-for="item in data.sourcesNavi" :key="item.name" :value="item.name" rounded="0">
          {{ item.title }}
        </v-tab>
      </v-tabs>
      <v-divider />
      <template v-if="source">
        <v-card-text style="position: relative;">
          <GraphCountPerNamespace v-if="showExpanded" :data="data.charts.namespaceScope[source].complete" />
          <GraphCountPerNamespace v-else :data="data.charts.namespaceScope[source].preview" />

          <v-btn v-if="hasPreview" variant="outlined" size="small" @click="expand = !expand" style="position: absolute; bottom: 10px; right: 10px;" rounded="0">
            <span v-if="showExpanded">Show preview</span>
            <span v-else>Show Complete List</span>
          </v-btn>
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
import { type Dashboard, Status, type Findings } from "~/modules/core/types";

const props = defineProps<{ data: Dashboard; hideCluster?: boolean }>();

const status1 = ref(Status.PASS);
const status2 = ref(Status.FAIL);

const source = ref('');

const expand = ref(false)

const hasPreview = computed(() => !!props.data.charts.namespaceScope[source.value].preview)

const showExpanded = computed(() => {
  if (!props.data.charts.namespaceScope[source.value].preview) {
    return true
  }

  return expand.value
})

watch(() => props.data.sources, (s) => {
  if (!s || !s.length || source.value) return

  source.value = s[0]
})
</script>
