<template>
  <v-container fluid class="py-4 px-4 main-height" :key="route.params.policy" v-if="data">
    <app-row>
      <v-card elevation="2">
        <v-toolbar color="header">
          <v-toolbar-title>
            {{ capilize(route.params.source) }}: {{ data.title }}
          </v-toolbar-title>
          <template #append>
            <v-btn variant="text" color="white" prepend-icon="mdi-close" @click="close">close</v-btn>
          </template>
        </v-toolbar>
      </v-card>
    </app-row>

    <policy-details :policy="data" v-if="data.showDetails" default-open />
    <app-row v-else>
      <v-card>
        <v-card-text>
          <v-alert type="error" variant="outlined">
            No additional Information available, ensure that you are running a required Plugin for the {{ route.params.source }} policies.
          </v-alert>
        </v-card-text>
      </v-card>
    </app-row>
  </v-container>
</template>

<script lang="ts" setup>
import { onChange } from "~/helper/compare";
import { useAPI } from "~/modules/core/composables/api";
import { capilize } from "~/modules/core/layouthHelper";

const route = useRoute()

const { data, refresh } = useAPI((api) => api.policyDetails(route.params.source, route.params.policy));

watch(route, onChange(refresh))

const close = () => window.close()
</script>
