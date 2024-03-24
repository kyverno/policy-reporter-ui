<template>
  <page-layout :title="`${capilize(route.params.source)}: ${data.title}`"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="route.params.source"
               v-if="data"
               hide-report
  >
    <template #append>
      <form-status-select style="min-width: 200px;" class="mr-2" />
      <v-btn variant="text" color="white" prepend-icon="mdi-arrow-left" @click="router.back()">back</v-btn>
    </template>

    <policy-details :policy="data" v-if="data.showDetails" />

    <policy-status-charts :data="data" :policy="route.params.policy" />
    <policy-cluster-results :source="route.params.source" :policy="route.params.policy" :status="status" />
    <policy-namespace-section
        :exceptions="data.exceptions"
        :namespaces="data.namespaces"
        :source="route.params.source"
        :policy="route.params.policy"
        :status="status"
    />
  </page-layout>
</template>

<script lang="ts" setup>
import { onChange } from "~/helper/compare";
import { useAPI } from "~/modules/core/composables/api";
import { capilize } from "~/modules/core/layouthHelper";
import { APIFilter } from "~/modules/core/provider/dashboard";

const route = useRoute()
const router = useRouter()

const { load } = useSourceStore(route.params.source)
await load(route.params.source)

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ kinds: [...kinds.value, ...clusterKinds.value] }))

provide(APIFilter, filter)

const status = computed(() => {
  if (!route.query.status) return undefined

  if (Array.isArray(route.query.status)) return route.query.status

  return [route.query.status]
})

const { data, refresh } = useAPI((api) => api.policyDetails(
    route.params.source,
    route.params.policy,
    route.query.namespace,
    status.value,
    filter.value.kinds,
));

watch(route, onChange(refresh))
watch(status, onChange(refresh))
watch(filter, onChange(refresh))
</script>
