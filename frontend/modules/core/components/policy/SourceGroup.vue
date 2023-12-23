<template>
  <v-row>
    <v-col cols="12">
      <v-card>
        <v-toolbar color="primary">
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
          <template  v-for="category in source.categories" :key="category.name">
            <v-divider />
            <policy-list :category="category.name" />
          </template>
        </div>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import { ResourceFilter } from "~/modules/core/provider/dashboard";
import type { Source } from "~/modules/core/types";
import CollapseBtn from "~/components/CollapseBtn.vue";

const props = defineProps<{ source: Source; }>();

const open = ref(true)

provide(ResourceFilter, ref({ sources: [props.source.name]}))

</script>
