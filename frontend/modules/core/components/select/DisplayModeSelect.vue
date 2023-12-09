<template>
  <v-select
    :model-value="(displayMode as Readonly<DisplayMode>)"
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
import { displayMode } from '~/store/session'

const input = (mode: DisplayMode) => { displayMode.value = mode }

</script>
