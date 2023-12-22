<template>
<v-container  v-if="data.counts.length" fluid class="py-4 px-4 main-height">
  <v-row>
    <v-col>
      <v-toolbar color="indigo" elevation="2" rounded>
        <v-toolbar-title>{{ capilize(route.params.source) }}: {{ route.params.category }}</v-toolbar-title>
        <template #append>
          <FormKindAutocomplete style="min-width: 300px; max-width: 100%; margin-right: 15px;" v-model="kinds" :source="route.params.source" />
          <FormClusterKindAutocomplete style="min-width: 300px;" v-model="clusterKinds" :source="route.params.source" />
        </template>
      </v-toolbar>
    </v-col>
  </v-row>
  <SourceStatus :data="data.counts[0]" :category="route.params.category" />
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
</v-container>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { clusterKinds, kinds } from '~/modules/core/store/filter';
import { capilize } from "~/modules/core/layouthHelper";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { ClusterKinds, NamespacedKinds, ResourceFilter } from "~/modules/core/provider/dashboard";

const route = useRoute()

const filter = ref({
  sources: [route.params.source],
  categories: [route.params.category],
})

const namespaces = await callAPI((api) => api.namespaces(route.params.source))

const { data, refresh } = useAPI((api) => api.countFindings({
      kinds: [...kinds.value, ...clusterKinds.value],
      ...filter.value,
    }),
    {
      default: () => ({ total: 0, counts: [] }),
    }
);

watch(kinds, () => refresh())

provide(ResourceFilter, ref(filter))
provide(NamespacedKinds, kinds)
provide(ClusterKinds, clusterKinds)
</script>
