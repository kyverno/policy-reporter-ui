<template>
  <page-layout v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               v-if="data"
  >
    <GraphSourceStatus v-if="data.singleSource" :data="data" :source="data.sources[0]" />
    <GraphSourcesStatus v-else :data="data" />
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :details="data.multiSource" />
      </v-col>
    </v-row>
    <resource-scroller :list="data.namespaces" v-if="data.namespaces.length">
      <template #default="{ item }">
        <LazyResourceResultList :namespace="item" :details="data.multiSource" />
      </template>
    </resource-scroller>
  </page-layout>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { onChange } from "~/helper/compare";
import { ResourceFilter } from "~/modules/core/provider/dashboard";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({
  kinds: kinds.value,
  clusterKinds: clusterKinds.value
}))

const { data, refresh } = useAPI((api) => api.dashboard(filter.value))

watch(filter, onChange(refresh))

provide(ResourceFilter, filter)
</script>
