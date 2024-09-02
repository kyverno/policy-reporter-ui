<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" :sm="Math.ceil(24 / showed.length)" :md="Math.ceil(12 / showed.length)" v-for="severity in showed" :key="severity">
        <v-card flat :title="severity" :class="['text-white', 'text-center', `bg-severity-${severity}`]">
          <v-card-text class="text-h3 my-4">
            {{ data[severity] }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import {Severity} from "../../types";

defineProps<{ data: { [key in Severity]: number; }; }>()

const severities = useSeveritiesInjection()

const showed = computed(() => severities.value)
</script>
