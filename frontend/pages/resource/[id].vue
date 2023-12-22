<template>
  <v-container fluid class="py-4 px-4 main-height" :key="route.params.id" v-if="data">
    <v-row>
      <v-col>
        <v-card elevation="2" rounded>
          <div  class="bg-indigo">
            <v-card-title>
              <v-container fluid class="ma-0 pa-0">
                <v-row>
                  <v-col>
                    <span v-if="data?.resource.namespace">{{ data?.resource.namespace }}/</span>{{ data.resource.name }}
                  </v-col>
                  <v-col class="text-right">
                    <v-btn variant="text" color="white" prepend-icon="mdi-arrow-left" @click="router.back()">back</v-btn>
                  </v-col>
                </v-row>
              </v-container>
            </v-card-title>
            <v-card-subtitle class="pb-4">{{ data?.resource.apiVersion }} {{ data.resource.kind }}</v-card-subtitle>
          </div>
          <v-card-text>
            <resource-result-counts :data="data.counts as any" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row v-if="data.sources.length > 1">
      <v-col>
        <v-card>
          <v-card-text>
            <resource-status :data="data.counts as any" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <resource-scroller :list="sources">
      <template #default="{ item }">
        <category-tables :source="item" :resource="data.resource" />
      </template>
    </resource-scroller>
  </v-container>
</template>

<script lang="ts" setup>
import ResourceStatus from "~/modules/core/components/chart/ResourceStatus.vue";
import ResourceResultCounts from "~/modules/core/components/chart/ResourceResultCounts.vue";
import type { Filter, Resource, ResourceStatusCount, Source } from "~/modules/core/types";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";

const route = useRoute()
const router = useRouter()

const filter = computed(() => {
  const f: Filter = {}
  if (route.query.source && typeof route.query.source === 'string') { f.sources = [route.query.source] }
  if (route.query.category && typeof route.query.category === 'string') { f.categories = [route.query.category] }

  if (route.query.sources && typeof route.query.source === 'object') { f.sources = route.query.sources }
  if (route.query.categories && typeof route.query.category === 'object') { f.categories = route.query.category }

  return f
})

const { data } = useAPI(
    async (api) => {
      let [resource, counts, sources] = await Promise.all([
        api.resource(route.params.id as string),
        api.resourceStatusCount(route.params.id as string, filter.value),
        api.sources(route.params.id as string)
      ]);

      if (route.query.source) {
        sources = sources.filter(s => s.name === route.query.source)
      }

      if (route.query.category) {
        sources = sources.map(s => ({
          name: s.name,
          categories: s.categories.filter(c => c.name === route.query.category)
        }))
      }

      return { resource, counts, sources: sources ?? [route.query.source] };
    }, {
      default: () => ({ resource: {} as Resource, counts: [] as ResourceStatusCount[], sources: [] as Source[] }),
    }
);

const sources = computed(() => (data.value?.sources || []).sort((a: Source, b: Source) => a.name.localeCompare(b.name)))
</script>
