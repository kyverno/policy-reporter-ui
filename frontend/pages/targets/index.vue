<template>
  <v-container fluid class="py-4 px-4 main-height" v-if="targets">
    <app-row>
      <v-card>
        <v-toolbar color="header" elevation="2">
          <v-toolbar-title>Configured Notification Targets</v-toolbar-title>
        </v-toolbar>
      </v-card>
    </app-row>
    <app-row v-for="(list, type) in targets" :key="type">
      <target-group :type="type" :targets="list" />
    </app-row>
  </v-container>
</template>

<script setup lang="ts">
const { data: targets } = useAPI((api) => api.targets())
</script>
