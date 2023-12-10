<template>
<v-card>
  <v-toolbar color="transparent">
    <v-toolbar-title>Cluster Resources</v-toolbar-title>
    <template #append>
      <Search class="mr-2" v-model="search" style="min-width: 300px;" />
      <SelectClusterKindAutocomplete style="min-width: 300px;" />
      <CollapseBtn v-model="open" />
    </template>
  </v-toolbar>
    <template  v-if="data?.items?.length && open">
      <v-list lines="two">
          <template v-for="item in data.items" :key="item.id">
              <v-divider />
              <v-list-item :to="`/resource/${item.id}`">
                  <v-list-item-title>{{ item.name }}</v-list-item-title>
                  <v-list-item-subtitle>{{ item.apiVersion }} {{ item.kind }}</v-list-item-subtitle>
                  <template v-slot:append>
                    <ResultChip :status="Status.PASS" :count="item.pass" tooltip="pass results" />
                    <ResultChip class="ml-2" :status="Status.WARN" :count="item.warn" tooltip="warning results" />
                    <ResultChip class="ml-2" :status="Status.FAIL" :count="item.fail" tooltip="fail results" />
                    <ResultChip class="ml-2" :status="Status.ERROR" :count="item.error" tooltip="error results" />
                  </template>
              </v-list-item>
          </template>
      </v-list>
      <template v-if="data.count > options.offset">
        <v-divider />
        <v-pagination v-model="options.page" :length="length" class="my-4" />
      </template>
    </template>
    <template v-if="!pending && !(data?.items?.length)">
        <v-divider />
        <v-card-text>
            No resources for the selected kinds found
        </v-card-text>
    </template>
</v-card>
</template>

<script setup lang="ts">
import { type Pagination, Status } from '../types'
import { clusterKinds } from "~/modules/core/store/filter";
import CollapseBtn from "~/components/CollapseBtn.vue";

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
    (api) => api.clusterResourceResults({ kinds: clusterKinds.value, search: search.value }, options),
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
