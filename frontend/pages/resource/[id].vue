<template>
  <v-container fluid class="py-4 px-4 main-height" :key="route.params.id" v-if="data">
    <v-row>
      <v-col>
        <v-card elevation="2" rounded>
          <div class="bg-header">
            <v-card-title>
              <v-container fluid class="ma-0 pa-0">
                <v-row>
                  <v-col>
                    <span v-if="data.resource.namespace">{{ data.resource.namespace }}/</span>{{ data.resource.name }}
                  </v-col>
                  <v-col class="text-right">
                    <v-btn variant="text" color="white" prepend-icon="mdi-arrow-left" @click="router.back()">back</v-btn>
                  </v-col>
                </v-row>
              </v-container>
            </v-card-title>
            <v-card-subtitle class="pb-4 text-grey-lighten-2" style="opacity: 1">{{ data.resource.apiVersion }} {{ data.resource.kind }}</v-card-subtitle>
          </div>
          <v-card-text>
            <resource-result-counts :data="data.results as any" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row v-if="data.chart">
      <v-col>
        <v-card>
          <v-card-text>
            <resource-status :data="data.chart" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <resource-scroller :list="data.sources">
      <template #default="{ item }">
        <category-tables :source="item" :resource="data.resource" />
      </template>
    </resource-scroller>
  </v-container>
</template>

<script lang="ts" setup>
import ResourceStatus from "~/modules/core/components/chart/ResourceStatus.vue";
import ResourceResultCounts from "~/modules/core/components/chart/ResourceResultCounts.vue";
import type { Filter } from "~/modules/core/types";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";
import { onChange } from "~/helper/compare";

const route = useRoute()
const router = useRouter()

const filter = computed(() => {
  const f: Filter = {}
  if (route.query.source && typeof route.query.source === 'string') { f.sources = [route.query.source] }
  if (route.query.category && typeof route.query.category === 'string') { f.categories = [route.query.category] }

  if (route.query.categories) {
    if (typeof route.query.categories === 'object') { f.categories = route.query.categories }
    if (typeof route.query.categories === 'string') { f.categories = [route.query.categories] }
  }

  if (route.query.sources) {
    if (typeof route.query.sources === 'object') { f.sources = route.query.sources }
    if (typeof route.query.sources === 'string') { f.sources = [route.query.sources] }
  }

  return f
})

const { data, refresh } = useAPI((api) => api.resource(route.params.id as string, filter.value));

watch(filter, onChange(refresh))
</script>
