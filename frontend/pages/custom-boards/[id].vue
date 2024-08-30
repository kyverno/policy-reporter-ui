<template>
  <page-layout v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :title="data.title"
               v-if="data"
               :source="data.singleSource ? data.sources[0] : undefined"
               :ns-scoped="!data.clusterScope"
               :store="id"
  >
    <template v-if="data.namespaces.length">
      <GraphSourceCharts :data="data" :hide-cluster="!data.clusterScope" />
      <v-row v-if="data.clusterScope">
        <v-col>
          <resource-cluster-list :details="data.multiSource" />
        </v-col>
      </v-row>
      <resource-namespace-section v-if="data.namespaces.length" :namespaces="data.namespaces">
        <template #default="{ namespaces }">
          <resource-scroller :list="namespaces" :default-loadings="3">
            <template #default="{ item }">
              <resource-list :namespace="item" :details="data.multiSource" />
            </template>
          </resource-scroller>
        </template>
      </resource-namespace-section>
    </template>
    <v-card class="mt-4" v-else>
      <v-card-text>
        <v-alert variant="outlined" type="error">
          No configured namespaces are found
        </v-alert>
      </v-card-text>
    </v-card>
  </page-layout>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { onChange } from "~/helper/compare";
import { APIFilter } from "~/modules/core/provider/dashboard";

const route = useRoute()

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({
  kinds: kinds.value,
  clusterKinds: clusterKinds.value
}))

const id = computed(() => route.params.id as string)

const { data, refresh } = useAPI((api) => api.customBoard(id.value, filter.value))

const store = useSourceStore(id.value)

watchEffect(() => {
  if (!data.value) return;

  store.load(data.value.sources)
})

watch(filter, onChange(refresh))

provide(APIFilter, computed(() => ({
  ...filter.value,
  sources: data.value?.filterSources,
})))

useStatusProvider(data)
useSeveritiesProvider(data)
useDashboardType(data)
</script>
