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

const props = defineProps<{ source: string; }>()

const filter = inject<Ref<Filter>>(ResourceFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(ClusterKinds, ref<string[]>([]))
const statusColors = useStatusColors()

const { data, refresh } = useAPI(
    (api) => api.statusCount(props.source, {
      ...filter.value,
      kinds: kinds.value.length ? kinds.value : undefined
    }), {
      default: () => ({
        [Status.PASS]: 0,
        [Status.WARN]: 0,
        [Status.FAIL]: 0,
        [Status.ERROR]: 0,
      }),
    }
);

const statusCounts = computed<{ [status in Partial<Status>]: number }>(() => ({
  [Status.PASS]: data.value.pass,
  [Status.WARN]: data.value.warn,
  [Status.FAIL]: data.value.fail,
  [Status.ERROR]: data.value.error,
}))


watch(kinds, () => refresh())
</script>
