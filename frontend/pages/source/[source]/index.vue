<template>
  <page-layout :title="capilize(route.params.source)" v-model:kinds="kinds" v-model:cluster-kinds="clusterKinds" :source="route.params.source">
    <SourceStatus :data="findings" />
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :source="route.params.source" :details="false" />
      </v-col>
    </v-row>
    <resource-scroller :list="namespaces">
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
import { execOnChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const route = useRoute()

const namespaces = await callAPI((api) => api.namespaces(route.params.source, {
  kinds: kinds.value.length ? kinds.value : undefined
}))

const filter = computed<Filter>(() => ({
  sources: [route.params.source],
  kinds: [...kinds.value, ...clusterKinds.value]
}))

const { data, refresh, pending } = useAPI(
    (api) => api.countFindings(filter.value),
    {
      default: () => ({ total: 0, counts: [{ source: route.params.source, counts: {}}] }),
    }
);

const findings = computed(() => {
  if (data.value?.counts?.length) return data.value.counts[0];

  return { source: route.params.source, counts: {}, total: 0 }
})

watch(filter, (n, o) => execOnChange(n, o, () => refresh()))

provide(ResourceFilter, filter)
</script>
