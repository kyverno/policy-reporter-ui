<template>
  <v-menu>
    <template #activator="{ props }">
      <v-btn variant="outlined" class="text-primary ml-3" rounded="4" height="40" v-bind="props">{{ profile.firstname }} {{ profile.lastname }}</v-btn>
    </template>
    <v-list class="pb-0">
      <v-list-item href="logout">
        <v-list-item-title>Logout</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import type { Profile } from "~/modules/core/types";

defineProps<{ profile: Profile }>()
</script>
