<template>
<v-container fluid class="py-4 px-4 main-height">
  <v-row>
    <v-col>
      <v-toolbar color="indigo" elevation="2" rounded>
        <v-toolbar-title>{{ capilize(route.params.source) }}: {{ route.params.category }}</v-toolbar-title>
        <template #append>
          <SelectKindAutocomplete style="width: 500px; max-width: 100%; margin-left: 15px;" />
        </template>
      </v-toolbar>
    </v-col>
  </v-row>
  <SourceStatus v-if="data.counts.length" :data="data.counts[0]" :category="route.params.category" />
  <v-row>
    <v-col>
      <LazyClusterResourceResultList :source="route.params.source" :category="route.params.category" />
    </v-col>
  </v-row>
  <v-infinite-scroll :onLoad="load" class="no-scrollbar">
    <template v-for="ns in loaded" :key="ns">
      <v-row>
        <v-col>
          <LazyResourceResultList :namespace="ns" :details="false" :source="route.params.source" :category="route.params.category" />
        </v-col>
      </v-row>
    </template>
    <template #empty></template>
  </v-infinite-scroll>
</v-container>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { clusterKinds, kinds } from '~/modules/core/store/filter';
import { capilize } from "~/modules/core/layouthHelper";

const route = useRoute()

const { data, refresh } = useAPI(
    (api) => api.countFindings({
      kinds: [...kinds.value, ...clusterKinds.value],
      sources: [route.params.source],
      categories: [route.params.category],
    }),
    {
      default: () => ({ total: 0, counts: [] }),
    }
);

const { data: namespaces } = useAPI(
    (api) => api.namespaces(route.params.source),
    {
      default: () => [],
    }
);

const loaded = ref<string[]>([])
const index = ref(1)

watch(namespaces, (ns: string[] | null) => {
  loaded.value = (ns || []).slice(0, 1)
})

const load = ({ done }: any) => {
  const sum = (namespaces.value || []).length

  const last = index.value
  const next = index.value + 2 > sum ? sum :  index.value + 2
  loaded.value = [...loaded.value, ...(namespaces.value || []).slice(last, next)]

  index.value = next
  if (next === sum) {
    done('empty')
  } else {
    done('ok')
  }
}


watch(kinds, refresh)
</script>
