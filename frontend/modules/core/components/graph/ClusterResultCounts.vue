<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="pass results" :class="['text-white', 'text-center', `bg-status-${Status.PASS}`]">
          <v-card-text class="text-h3 my-4">
            {{ data.pass }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="warn results" :class="['text-white', 'text-center', `bg-status-${Status.WARN}`]">
          <v-card-text class="text-h3 my-4">
            {{ data.warn }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="fail results" :class="['text-white', 'text-center', `bg-status-${Status.FAIL}`]">
          <v-card-text class="text-h3 my-4">
            {{ data.fail }}
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" sm="6" md="3">
        <v-card flat title="error results" :class="['text-white', 'text-center', `bg-status-${Status.ERROR}`]">
          <v-card-text class="text-h3 my-4">
            {{ data.error }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { Status } from "../../types";

const props = defineProps<{ data: { [key in Status]: number; }; }>()
</script>
