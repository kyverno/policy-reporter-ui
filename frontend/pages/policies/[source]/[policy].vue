<template>
  <v-container fluid class="py-4 px-4 main-height" :key="route.params.policy" v-if="data">
    <app-row>
      <v-card elevation="2">
        <v-toolbar color="header">
          <v-toolbar-title>
            {{ capilize(route.params.source) }}: {{ data.title }}
          </v-toolbar-title>
          <template #append>
            <form-status-select style="min-width: 200px;" class="mr-2" />
            <v-btn variant="text" color="white" prepend-icon="mdi-arrow-left" @click="router.back()">back</v-btn>
          </template>
        </v-toolbar>
      </v-card>
    </app-row>

    <policy-details :policy="data" v-if="data.showDetails" />

    <policy-status-charts :data="data" :policy="route.params.policy" />
    <policy-cluster-results :source="route.params.source" :policy="route.params.policy" :status="status" />
    <policy-namespace-section :namespaces="data.namespaces" :source="route.params.source" :policy="route.params.policy" :status="status" />
  </v-container>
</template>

<script lang="ts" setup>
import { onChange } from "~/helper/compare";
import { useAPI } from "~/modules/core/composables/api";
import { capilize } from "~/modules/core/layouthHelper";

const route = useRoute()
const router = useRouter()

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
));

watch(route, onChange(refresh))
watch(status, onChange(refresh))
</script>
