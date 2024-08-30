<template>
  <v-row>
    <v-col cols="12" md="6">
      <v-card>
        <v-card-text class="mb-0 pb-0">
          <GraphSeverityFindings :data="(data.charts.findings[severity1] as Findings)" :severity="severity1" :key="severity1" />
        </v-card-text>
        <v-card-actions style="margin-top: -40px">
          <v-btn-toggle
              v-model="severity1"
              rounded="0"
              mandatory
              style="height: 40px"
              divided
              variant="outlined"
          >
            <severity-btn v-if="data.total.perResult[Severity.MEDIUM]" :severity="Severity.MEDIUM" />
            <severity-btn v-if="data.total.perResult[Severity.LOW]" :severity="Severity.LOW" />
            <severity-btn v-if="data.total.perResult[Severity.INFO]" :severity="Severity.INFO" />
          </v-btn-toggle>
        </v-card-actions>
      </v-card>
    </v-col>
    <v-col cols="12" md="6">
      <v-card>
        <v-card-text class="mb-0 pb-0">
          <GraphSeverityFindings :data="(data.charts.findings[severity2] as Findings)" :severity="severity2" :key="severity2" :time="1200" />
        </v-card-text>
        <v-card-actions style="margin-top: -40px">
          <v-btn-toggle
              v-model="severity2"
              rounded="0"
              mandatory
              style="height: 40px"
              divided
              variant="outlined"
          >
            <severity-btn v-if="data.total.perResult[Severity.HIGH]" :severity="Severity.HIGH" />
            <severity-btn v-if="data.total.perResult[Severity.CRITICAL]" :severity="Severity.CRITICAL" />
            <severity-btn v-if="data.total.perResult[Severity.UNKNOWN]" :severity="Severity.UNKNOWN" />
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
            <GraphClusterSeverityCounts :data="data.charts.clusterScope[source]" class="px-0 pb-0" />
          </v-card-text>
        </template>
      </template>
    </v-card>
  </app-row>
</template>

<script setup lang="ts">
import { type Dashboard, type Findings, Severity } from "~/modules/core/types";

const props = defineProps<{ data: Dashboard; hideCluster?: boolean }>();

const severity1 = ref(Severity.MEDIUM);
const severity2 = ref(Severity.HIGH);

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
