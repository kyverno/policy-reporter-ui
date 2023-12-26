<template>
  <v-container fluid>
    <v-row>
      <v-col v-for="(count, status) in statusCounts" :key="status" cols="12" sm="6" md="3">
        <v-card flat :title="`${status} results`" class="text-white text-center" :style="`background-color: ${statusColors[status]}`">
          <v-card-text class="text-h3 my-4">
            {{ count }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { type Filter, Status } from "~/modules/core/types";
import type { Ref } from "vue";
import { ClusterKinds, ResourceFilter } from "~/modules/core/provider/dashboard";
import { useStatusColors } from "~/modules/core/composables/theme";

const props = defineProps<{ source?: string; }>()

const filter = inject<Ref<Filter>>(ResourceFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(ClusterKinds, ref<string[]>([]))
const statusColors = useStatusColors()

const { data: sc, refresh } = useAPI(
    (api) => api.statusCount({
      ...filter.value,
      sources: props.source ? [props.source] : undefined,
      kinds: kinds.value.length ? kinds.value : undefined
    }), {
      default: () => [
        { status: Status.PASS, count: 0 },
        { status: Status.WARN, count: 0 },
        { status: Status.FAIL, count: 0 },
        { status: Status.ERROR, count: 0 },
      ],
    }
);

const statusCounts = computed<{ [status in Status]: number }>(() => {
  return sc.value.reduce((acc, item) => {
    if (item.status === Status.SKIP) return acc;

    acc[item.status] = item.count

    return acc
  }, {
    [Status.PASS]: 0,
    [Status.WARN]: 0,
    [Status.FAIL]: 0,
    [Status.ERROR]: 0,
  })
})


watch(kinds, () => refresh())
</script>
