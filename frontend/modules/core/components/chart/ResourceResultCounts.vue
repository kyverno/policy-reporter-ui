<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="pass results" class="text-white text-center" :style="`background-color: ${statusColors.pass}`">
          <v-card-text class="text-h3 my-4">
            {{ data[Status.PASS] }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="warn results" class="text-white text-center" :style="`background-color: ${statusColors.warn}`">
          <v-card-text class="text-h3 my-4">
            {{ data[Status.WARN] }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="fail results" class="text-white text-center" :style="`background-color: ${statusColors.fail}`">
          <v-card-text class="text-h3 my-4">
            {{ data[Status.FAIL] }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="error results" class="text-white text-center" :style="`background-color: ${statusColors.error}`">
          <v-card-text class="text-h3 my-4">
            {{ data[Status.ERROR] }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { Status } from "~/modules/core/types";
import { useStatusColors } from "~/modules/core/composables/theme";

const props = defineProps<{ data: { [key in Status]: number; }}>()

const statusColors = useStatusColors()
</script>
