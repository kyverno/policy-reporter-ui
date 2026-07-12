<template>
  <page-layout :title="`${capilize(source)}: ${data.title}`"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="source"
               v-if="data"
               hide-report
  >
    <template #append>
      <form-status-select style="min-width: 200px;" class="mr-2" />
      <v-btn variant="text" color="white" prepend-icon="mdi-arrow-left" @click="router.back()">back</v-btn>
    </template>

    <policy-details :policy="data" v-if="data.showDetails" />

    <policy-status-charts :data="data" :policy="policy" />
    <policy-cluster-results :source="source" :policy="policy" :status="status" />
    <policy-namespace-section
        :exceptions="data.exceptions"
        :namespaces="data.namespaces"
        :source="source"
        :policy="policy"
        :status="status"
    />
  </page-layout>
</template>

<script lang="ts" setup>
import { APIFilter } from "~/provider/dashboard";
import { Status } from "~/types/core";

const route = useRoute()
const router = useRouter()

const { filter, clusterKinds, kinds, source, policy } = useFilter()

const { load } = useSourceStore(source.value)
await load(source.value)


provide(APIFilter, filter)

const status = computed<Status[] | undefined>(() => {
  if (!route.query.status) return undefined

  if (Array.isArray(route.query.status)) return route.query.status as Status[]

  return [route.query.status] as Status[]
})

const { data, refresh } = useAPI((api) => api.policyDetails(
    source.value,
    policy.value,
    route.query.namespace as string,
    status.value,
    filter.value.kinds,
));

watch(route, onChange(refresh))
watch(status, onChange(refresh))
watch(filter, onChange(refresh))

useStatusProvider(data)
</script>
