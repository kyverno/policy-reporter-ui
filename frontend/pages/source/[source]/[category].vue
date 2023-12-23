<template>
  <page-layout v-if="data.counts.length"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="route.params.source"
               :title="`${capilize(route.params.source)}: ${ route.params.category }`"
  >
    <SourceStatus :category="route.params.category" :data="data.counts[0]"/>
    <v-row>
      <v-col>
        <LazyClusterResourceResultList :details="false" :source="route.params.source"/>
      </v-col>
    </v-row>
    <resource-scroller :list="namespaces">
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
import { execOnChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const route = useRoute()

const filter = computed(() => ({
  sources: [route.params.source],
  categories: [route.params.category],
  kinds: [...kinds.value, ...clusterKinds.value],
}))

const namespaces = await callAPI((api) => api.namespaces(route.params.source))

const { data, refresh } = useAPI((api) => api.countFindings(filter.value),
    {
      default: () => ({ total: 0, counts: [] }),
    }
);

watch(filter, (n, o) => execOnChange(n, o, () => refresh()))

provide(ResourceFilter, ref(filter))
</script>
