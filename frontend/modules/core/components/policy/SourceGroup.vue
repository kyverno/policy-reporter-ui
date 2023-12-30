<template>
  <v-row>
    <v-col cols="12">
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
              <policy-list :category="item" />
            </template>
          </scroller>
        </div>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ResourceFilter } from "~/modules/core/provider/dashboard";
import type { SourceDetails } from "~/modules/core/types";
import CollapseBtn from "~/components/CollapseBtn.vue";
import ResourceScroller from "~/modules/core/components/ResourceScroller.vue";

const props = defineProps<{ source: SourceDetails; }>();

const open = ref(true)

provide(ResourceFilter, ref({ sources: [props.source.name]}))

</script>
