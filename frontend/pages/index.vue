<template>
  <v-container fluid v-if="data" class="py-4 px-4 main-height">
    <v-row>
      <v-col>
        <v-toolbar color="indigo" elevation="2" rounded>
          <template #append>
            <SelectKindAutocomplete style="width: 500px; max-width: 100%; margin-left: 15px;" />
          </template>
        </v-toolbar>
      </v-col>
    </v-row>
    <SourcesStatus v-if="multiSource" :data="data as FindingCounts" />
    <SourceStatus v-if="singleSource && data.counts.length > 0" :data="data.counts[0]" />
    <v-row>
      <v-col>
        <LazyClusterResourceResultList />
      </v-col>
    </v-row>
    <v-infinite-scroll :onLoad="load" class="no-scrollbar">
      <template v-for="ns in loaded" :key="ns">
        <v-row>
          <v-col>
            <LazyResourceResultList :namespace="ns" :details="multiSource" />
          </v-col>
        </v-row>
      </template>
      <template #empty></template>
    </v-infinite-scroll>
  </v-container>
</template>

<script setup lang="ts">
import { useAPI } from '~/modules/core/composables/api'
import { kinds } from '~/modules/core/store/filter';
import { type FindingCounts } from "~/modules/core/types";

const props = defineProps({
  sources: { types: Array, default: () => [] }
})

const { data, refresh } = useAPI(
    (api) => api.countFindings({ kinds: kinds.value }),
    {
      default: () => ({ total: 0, counts: [] }),
    }
);


const multiSource = computed(() => (props.sources.length || 0) > 1)
const singleSource = computed(() => (props.sources.length || 0) === 1)

const { data: namespaces } = useAPI(
    (api) => api.namespaces(),
    {
      default: () => [],
    }
);

const loaded = ref<string[]>([])
const index = ref(0)

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
