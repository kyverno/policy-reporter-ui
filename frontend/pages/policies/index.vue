<template>
  <page-layout v-if="sources?.length"
               title="Policy Dashboard"
               v-model:kinds="kinds"
               v-model:cluster-kinds="clusterKinds"
  >
    <resource-scroller :list="sources">
      <template #default="{ item }">
        <policy-source-group :source="item" />
      </template>
    </resource-scroller>
  </page-layout>
</template>

<script setup lang="ts">
import { callAPI } from "~/modules/core/composables/api";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { ClusterKinds, NamespacedKinds } from "~/modules/core/provider/dashboard";
import { execOnChange } from "~/helper/compare";

const kinds = ref<string[]>([])
const clusterKinds = ref<string[]>([])

const filter = computed(() => ({ kinds: [...kinds.value, ...clusterKinds.value] }))

const { data: sources, refresh } = useAPI(
    (api) => api.sources(undefined, filter.value).then(s => s.sort((a, b) => a.name.localeCompare(b.name))),
    { default: () => [] }
)

watch(filter, (a, b) => execOnChange(a, b, () => refresh()))
</script>
