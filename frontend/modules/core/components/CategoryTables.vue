<template>
  <v-row>
    <v-col cols="12">
      <v-card>
        <v-toolbar color="category">
          <v-toolbar-title>{{ capilize(source.name) }}</v-toolbar-title>
          <template #append>
            <CollapseBtn v-model="open" />
          </template>
        </v-toolbar>
        <div v-show="open">
          <v-divider />
          <v-card-text>
            <ChartStatusPerCategory :source="source" />
          </v-card-text>
          <result-table :source="source.name" :resource="resource.id" v-for="category in source.categories" :key="category.name" :category="category.name" />
        </div>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import type { Resource, Source } from "~/modules/core/types";
import CollapseBtn from "~/components/CollapseBtn.vue";

const open = ref(true)

const props = defineProps<{
  source: Source;
  resource: Resource;
}>();

</script>
