<template>
<v-card>
  <v-toolbar color="transparent">
    <v-toolbar-title>Cluster Resources</v-toolbar-title>
    <template #append>
      <Search class="mr-2" v-model="search" style="min-width: 300px;" />
      <SelectClusterKindAutocomplete style="min-width: 300px;" :source="props.source" />
      <CollapseBtn v-model="open" />
    </template>
  </v-toolbar>
  <v-list v-if="data?.items?.length && open" lines="two">
    <ResourceResultItem v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" />
  </v-list>
  <template v-if="!pending && !(data?.items?.length)">
      <v-divider />
      <v-card-text>
          No resources for the selected kinds found
      </v-card-text>
  </template>
</v-card>
</template>

<script setup lang="ts">
import { type Filter, type Pagination, Status } from '../types'
import { clusterKinds } from "~/modules/core/store/filter";
import CollapseBtn from "~/components/CollapseBtn.vue";

const props = defineProps<{ source?: string; filter?: Filter; details: boolean }>()

const search = ref('')
const open = ref(true)

const options = reactive<Pagination>({
  page: 1,
  offset: 15,
})

const length = computed(() => {
  return Math.ceil((data.value?.count || 0) / options.offset)
})

const { data, refresh, pending } = useAPI(
    (api) => api.clusterResourceResults({
      ...(props.filter || {}),
      kinds: clusterKinds.value,
      search: search.value,
    }, options),
    {
      default: () => ({ items: [], count: 0 }),
    }
);

watch(clusterKinds, () => {
  if (options.page !== 1) {
    options.page = 1
    return
  }

  refresh()
})

watch(search, () => {
  if (options.page !== 1) {
    options.page = 1
    return
  }

  refresh()
})

watch(options, refresh)

</script>
