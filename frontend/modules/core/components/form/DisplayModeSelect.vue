<template>
  <v-select
    :model-value="config.theme"
    :items="items"
    variant="outlined"
    hide-details
    density="compact"
    prepend-inner-icon="mdi-brightness-6"
    style="width: 220px;"
    @update:model-value="input"
  />
</template>

<script lang="ts" setup>
import { DisplayMode } from '~/modules/core/types'
import { useConfigStore } from "~/store/config";

const config = useConfigStore()

const items = [
  { title: 'Light', value: DisplayMode.LIGHT },
  { title: 'Dark', value: DisplayMode.DARK },
  { title: 'Color Blind Light', value: DisplayMode.COLOR_BLIND },
  { title: 'Color Blind Dark', value: DisplayMode.COLOR_BLIND_DARK },
]

const input = (mode: DisplayMode) => { config.setDisplayMode(mode) }

</script>
