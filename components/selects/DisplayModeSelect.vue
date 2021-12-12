<template>
  <v-select
    :value="displayMode"
    :items="modes"
    outlined
    hide-details
    dense
    prepend-inner-icon="mdi-brightness-6"
    style="max-width: 140px;"
    @input="setDisplayMode"
  />
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters, mapMutations } from 'vuex'
import { DisplayMode } from '~/policy-reporter-plugins/core/types'

type Data = {
  modes: DisplayMode[];
}

type Methdos = {
  setDisplayMode(mode: DisplayMode): void;
}

type Computed = {
  displayMode: DisplayMode;
}

export default Vue.extend<Data, Methdos, Computed, {}>({
  data: () => ({
    modes: [DisplayMode.LIGHT, DisplayMode.DARK]
  }),
  computed: mapGetters(['displayMode']),
  methods: mapMutations(['setDisplayMode'])
})
</script>
