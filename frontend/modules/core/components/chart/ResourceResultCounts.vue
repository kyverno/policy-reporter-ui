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
import { type ResourceStatusCount, Status } from "~/modules/core/types";
import { useStatusColors } from "~/modules/core/composables/theme";

const props = defineProps<{ data: ResourceStatusCount[] }>()

const statusColors = useStatusColors()
const statusCounts = computed(() => {
  return props.data?.reduce<{ [status in Omit<Status, 'skip'>]: number }>((acc, item) => {
    if (item.status === Status.SKIP) return acc;

    acc[item.status] = item.items.reduce((s, i) => s + i.count, 0)

    return acc
  }, {
    [Status.PASS]: 0,
    [Status.WARN]: 0,
    [Status.FAIL]: 0,
    [Status.ERROR]: 0,
  })
})
</script>
