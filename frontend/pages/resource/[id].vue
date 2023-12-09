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
    <category-tables v-for="source in sources" :source="source" :resource="data.resource" :key="source" />
  </v-container>
</template>

<script lang="ts" setup>
import ResourceStatus from "~/modules/core/components/chart/ResourceStatus.vue";
import ResourceResultCounts from "~/modules/core/components/chart/ResourceResultCounts.vue";
import type { Source } from "~/modules/core/types";

const route = useRoute()
const router = useRouter()

const { data } = useAPI(
    (api) => Promise.all([
      api.resource(route.params.id as string),
      api.resourceStatusCount(route.params.id as string),
      api.sources(route.params.id as string),
    ]).then(([resource, counts, sources]) => ({ resource, counts, sources })), {
      default: () => ({ resource: {}, counts: [], sources: [] }),
    }
);

const sources = computed(() => (data.value?.sources || []).sort((a: Source, b: Source) => a.name.localeCompare(b.name)))
</script>
