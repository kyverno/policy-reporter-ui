<template>
  <v-container fluid class="py-4 px-4 main-height" :key="policy" v-if="data">
    <app-row>
      <v-card elevation="2">
        <v-toolbar color="header">
          <v-toolbar-title>
            {{ capilize(source) }}: {{ data.title }}
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
            No additional Information available, ensure that you are running a required Plugin for the {{ source }} policies.
          </v-alert>
        </v-card-text>
      </v-card>
    </app-row>
  </v-container>
</template>

<script lang="ts" setup>
const route = useRoute()
const source = computed<string>(() => route.params.source as string)
const policy = computed<string>(() => route.params.policy as string)

const { data, refresh } = useAPI((api) => api.policyDetails(source.value, policy.value));

watch(route, onChange(refresh))

const close = () => window.close()
</script>
