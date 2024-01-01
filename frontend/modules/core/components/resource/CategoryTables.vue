<template>
  <app-row>
    <v-card>
      <v-toolbar color="category">
        <v-toolbar-title>{{ source.title }}</v-toolbar-title>
        <template #append>
          <CollapseBtn v-model="open" />
        </template>
      </v-toolbar>
      <div v-show="open">
        <v-divider />
        <v-card-text>
          <GraphStatusPerCategory :source="source.chart" />
        </v-card-text>
        <scroller :list="source.categories">
          <template #default="{ item }">
            <resource-results :source="source.name" :resource="resource.id" :category="item" />
          </template>
        </scroller>
      </div>
    </v-card>
  </app-row>
</template>

<script setup lang="ts">
import type { Resource, SourceDetails } from "../../types";
import CollapseBtn from "../../../../components/CollapseBtn.vue";

const open = ref(true)

const props = defineProps<{
  source: SourceDetails;
  resource: Resource;
}>();

</script>
