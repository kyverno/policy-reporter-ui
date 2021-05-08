<template>
    <router-view />
</template>

<script lang="ts">
import Vue from 'vue';
import { FETCH_POLICIES, NAMESPACE } from '../store';

export default Vue.extend({
  data: () => ({
    interval: 0,
  }),
  created() {
    this.$store.dispatch(`${NAMESPACE}/${FETCH_POLICIES}`);

    // @ts-ignore
    this.interval = setInterval(() => {
      this.$store.dispatch(`${NAMESPACE}/${FETCH_POLICIES}`);
    }, 5000);
  },
  destroyed() {
    clearInterval(this.interval);
  },
});
</script>
