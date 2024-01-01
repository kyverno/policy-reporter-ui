<template>
  <v-row>
    <v-col>
      <slot />
    </v-col>
  </v-row>
</template>
<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
</script>
