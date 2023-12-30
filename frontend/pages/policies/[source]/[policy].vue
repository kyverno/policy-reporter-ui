<template>
  <v-container fluid class="py-4 px-4 main-height" :key="route.params.id" v-if="data">
    <v-row>
      <v-col>
        <v-card elevation="2" rounded>
          <v-toolbar color="header">
            <v-toolbar-title>
              {{ route.params.source }}: {{ route.params.policy }}
            </v-toolbar-title>
            <template #append>
              <v-btn variant="text" color="white" prepend-icon="mdi-arrow-left" @click="router.back()">back</v-btn>
            </template>
          </v-toolbar>
        </v-card>
      </v-col>
    </v-row>
    <policy-status-charts :data="data" :policy="route.params.policy" />
    <policy-cluster-results :source="route.params.source" :policy="route.params.policy" />
    <resource-scroller :list="data.namespaces">
      <template #default="{ item }">
        <v-row>
          <v-col>
            <policy-results :namespace="item" :source="route.params.source" :policy="route.params.policy" />
          </v-col>
        </v-row>
      </template>
    </resource-scroller>
  </v-container>
</template>

<script lang="ts" setup>
import { onChange } from "~/helper/compare";
import { useAPI } from "~/modules/core/composables/api";

const route = useRoute()
const router = useRouter()

const { data, refresh } = useAPI((api) => api.policyDetails(route.params.source, route.params.policy));

watch(route, onChange(refresh))
</script>
