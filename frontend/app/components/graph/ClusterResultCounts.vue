<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" :md="Math.ceil(12 / showed.length)" v-for="status in showed" :key="status">
        <v-card flat :title="status" :class="['text-white', 'text-center', `bg-status-${status}`]">
          <v-card-text class="text-display-large my-4">
            {{ data[status] }}
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { Status } from "~/types/core";

defineProps<{ data: { [key in Status]: number; }; }>()

const status = useStatusInjection()

const showed = computed(() => status.value.filter(s => s !== Status.SKIP))
</script>
