<template>
  <page-layout v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               v-model:mode="mode"
               :title="data.title"
               v-if="data"
               :source="source"
               :ns-scoped="!data.clusterScope"
               :store="id"
  >
    <template v-if="data.namespaces.length">
      <template v-if="isCompact">
        <LazyGraphSourceCard
          :data="data.charts.clusterScope[source!]!"
          :source="source!"
          :hide-cluster="!data.clusterScope"
          :type="dataType"
        />
        <resource-summary-list :details="false" class="mt-6" :data="data.summary" :custom-board="id" />
      </template>
      <template v-else>
        <GraphSourceCharts :data="data" :hide-cluster="!data.clusterScope" />
        <v-row v-if="data.clusterScope">
          <v-col>
            <custom-board-cluster-table v-if="showResults" :sources="data.sources" :id="id" />
            <custom-board-cluster-list v-else :details="data.multiSource" :id="id" />
          </v-col>
        </v-row>
        <resource-namespace-section v-if="data.namespaces.length" :namespaces="data.namespaces">
          <template #default="{ namespaces }">
            <resource-scroller :list="namespaces" :default-loadings="3">
              <template #default="{ item }">
                <custom-board-table v-if="showResults" :namespace="item" :sources="data.sources" :id="id" />
                <custom-board-list v-else :namespace="item" :details="data.multiSource" :id="id" />
              </template>
            </resource-scroller>
          </template>
        </resource-namespace-section>
      </template>
    </template>
    <v-card class="mt-4" v-else>
      <v-card-text>
        <v-alert variant="outlined" type="error">
          No configured namespaces are found
        </v-alert>
      </v-card-text>
    </v-card>
  </page-layout>
  <unauthorized v-if="error?.status === 401" />
</template>

<script setup lang="ts">
import { APIFilter } from "~/provider/dashboard";

const route = useRoute()

const { kinds, clusterKinds, filter } = useFilter()

const id = computed(() => route.params.id as string)

const { data, refresh, error } = useAPI((api) => api.customBoard(id.value, filter.value))
const { showResults, dataType, mode, isCompact } = useDashboardHelper(data)

const source = computed(() => data.value.singleSource ? data.value.sources[0] : undefined)

watch(filter, onChange(refresh))

provide(APIFilter, computed(() => ({
  ...filter.value,
  sources: data.value?.sources,
  namespaces: data.value?.namespaces,
  apis: data.value?.filter.resources,
  clusterApis: data.value?.filter.clusterResources,
})))

useDashboardProvider(data)
</script>
