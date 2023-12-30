<template>
  <page-layout v-if="data"
               v-model:cluster-kinds="clusterKinds"
               v-model:kinds="kinds"
               :source="route.params.source"
               :title="`${capilize(route.params.source)}: ${ route.params.category }`"
  >
    <GraphSourceStatus :category="route.params.category" :data="data" :source="route.params.source"/>
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :details="false" :source="route.params.source"/>
      </v-col>
    </v-row>
    <resource-scroller :list="data.namespaces">
      <template #default="{ item }">
        <LazyResourceResultList :details="false" :namespace="item"/>
      </template>
    </resource-scroller>
  </page-layout>
</template>

<script lang="ts" setup>
import { useAPI } from '~/modules/core/composables/api'
import { capilize } from "~/modules/core/layouthHelper";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { ResourceFilter } from "~/modules/core/provider/dashboard";
import { onChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const route = useRoute()

const filter = computed(() => ({
  sources: [route.params.source],
  categories: [route.params.category],
  kinds: kinds.value,
  clusterKinds: clusterKinds.value,
}))

const { data, refresh } = useAPI((api) => api.dashboard(filter.value));

watch(filter, onChange(refresh))

provide(ResourceFilter, ref(filter))
</script>
