<template>
  <v-card height="300px">
    <v-card-title class="pb-0">
      {{ statusText }} ClusterPolicies
    </v-card-title>
    <v-card-text style="height: calc(100% - 48px); font-size: 11rem!important;"
                :class="['text-center text-h1 d-flex justify-center align-center', `${color}--text`]"
    >
      {{ count }}
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend<{}, {}, {}, { count: number; statusText: string; color: string }>({
  name: 'FailingClusterPolicies',
  props: {
    statusText: { required: true, type: String },
    color: { default: 'default', type: String },
    count: { required: true, type: Number },
  },
});
</script>
