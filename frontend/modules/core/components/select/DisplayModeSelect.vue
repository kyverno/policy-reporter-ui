<template>
  <v-select
    :model-value="config.theme"
    :items="[DisplayMode.LIGHT, DisplayMode.DARK]"
    variant="outlined"
    hide-details
    density="compact"
    prepend-inner-icon="mdi-brightness-6"
    style="max-width: 140px;"
    @update:model-value="input"
  />
</template>

<script lang="ts" setup>
import { DisplayMode } from '../../types'
import { useConfigStore } from "~/store/config";

const config = useConfigStore()

const input = (mode: DisplayMode) => { config.setDisplayMode(mode) }

</script>
