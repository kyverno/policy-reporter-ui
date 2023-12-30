<template>
  <page-layout :title="capilize(route.params.source)" v-model:kinds="kinds" v-model:cluster-kinds="clusterKinds" :source="route.params.source" v-if="data">
    <GraphSourceStatus :data="data" :source="route.params.source" />
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :source="route.params.source" :details="false" />
      </v-col>
    </v-row>
    <resource-scroller :list="data.namespaces">
      <template #default="{ item }">
        <LazyResourceResultList :namespace="item" :details="false" />
      </template>
    </resource-scroller>
  </page-layout>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { capilize } from "~/modules/core/layouthHelper";
import { type Filter } from "~/modules/core/types";
import { ResourceFilter } from "~/modules/core/provider/dashboard";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { onChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const route = useRoute()

const filter = computed<Filter>(() => ({
  sources: [route.params.source],
  kinds: kinds.value,
  clusterKinds: clusterKinds.value,
}))

const { data, refresh } = useAPI((api) => api.dashboard(filter.value))

watch(filter, onChange(refresh))

provide(ResourceFilter, filter)
</script>
