<template>
  <page-layout v-if="sources"
               :title="capilize(route.params.source)"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
               :source="route.params.source"
  >
    <resource-scroller :list="sources">
      <template #default="{ item }">
        <policy-source-group :source="item" />
      </template>
    </resource-scroller>
  </page-layout>
  <unauthorized v-if="error?.status === 401" />
</template>

<script setup lang="ts">
import { onChange } from "~/helper/compare";
import { capilize } from "~/modules/core/layouthHelper";

const route = useRoute()

const { load } = useSourceStore(route.params.source as string)
await load(route.params.source)

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ sources: [route.params.source as string], kinds: [...kinds.value, ...clusterKinds.value] }))

const { data: sources, refresh, error } = useAPI((api) => api.policySources(filter.value))

watch(filter, onChange(refresh))
</script>
