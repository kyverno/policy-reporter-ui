<template>
  <v-select
    :value="refreshInterval"
    :items="intervals"
    outlined
    hide-details
    dense
    prepend-inner-icon="mdi-clock-outline"
    style="max-width: 120px;"
    @input="setRefreshInterval"
  />
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters, mapMutations } from 'vuex'

type Data = {
  intervals: Array<{ text: string; value: number }>
}

type Methdos = {
  setRefreshInterval(): void;
}

type Computed = {
  refreshInterval: number;
}

export default Vue.extend<Data, Methdos, Computed, {}>({
  data: () => ({
    intervals: [
      { text: '1s', value: 1000 },
      { text: '5s', value: 5000 },
      { text: '10s', value: 10000 },
      { text: '30s', value: 30000 },
      { text: '1min', value: 100000 }
    ]
  }),
  computed: mapGetters(['refreshInterval']),
  methods: mapMutations(['setRefreshInterval'])
})
</script>
