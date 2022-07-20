<template>
  <div>
    <slot v-if="ready" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend<{ ready: boolean }, {}, {}, { time: number }>({
  props: {
    time: { type: Number, default: 100 }
  },
  data: () => ({
    ready: false
  }),
  created () {
    setTimeout(() => { this.ready = true }, this.time)
  }
})
</script>
