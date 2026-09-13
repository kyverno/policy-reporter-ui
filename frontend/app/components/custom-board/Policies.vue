<template>
  <div v-if="!pending && !sources.length" class="mt-4">
    <v-card>
      <v-card-text>No policies found for this custom board.</v-card-text>
    </v-card>
  </div>

  <policy-source-group
    v-for="source in sources"
    :key="source.name"
    :source="source"
  />
</template>

<script setup lang="ts">
import { APIFilter } from "~/provider/dashboard";
import type { Filter, PolicyFilter, SourceDetails } from "~/types/core";
import type { Ref } from "vue";

const props = defineProps<{ id: string }>()

const filter = inject<Ref<Filter>>(APIFilter, ref<Filter>({}))

const { data, refresh, pending } = useAPI(
  (api) => api.customBoardPolicySources(props.id, filter.value),
  {
    default: (): { filter: PolicyFilter; sources: SourceDetails[] } => ({
      filter: {
        status: [],
        severities: [],
        namespaceKinds: [],
        clusterKinds: [],
      },
      sources: [],
    }),
  },
)

const sources = computed(() => data.value?.sources || [])

watch(filter, onChange(refresh))
</script>
