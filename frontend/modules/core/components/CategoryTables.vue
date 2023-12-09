<template>
  <v-row>
    <v-col cols="12">
      <v-card>
        <v-card-title class="py-3">
          <span class="text-h4">{{ capilize(source.name) }}</span>
        </v-card-title>
      </v-card>
    </v-col>
  </v-row>
  <v-row>
    <v-col cols="12">
      <v-card>
        <v-card-text>
          <ChartStatusPerCategory :source="source" />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
  <result-table :source="source.name" :resource="resource.id" v-for="category in source.categories" :key="category.name" :category="category.name" />
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import type { Resource, Source } from "~/modules/core/types";

const props = defineProps<{
  source: Source;
  resource: Resource;
}>();

</script>
