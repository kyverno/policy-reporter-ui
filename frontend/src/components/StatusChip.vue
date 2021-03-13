<template>
    <v-chip :color="color" :dark="dark" v-bind="$attrs" v-on="$listeners">{{ status }}</v-chip>
</template>

<script lang="ts">
import { Status } from '@/models';
import Vue, { PropType } from 'vue';

export default Vue.extend<{}, {}, { color: string; dark: boolean }, { status: Status }>({
  props: {
    status: { type: String as PropType<Status>, required: true },
  },
  inheritAttrs: false,
  computed: {
    color() {
      switch (this.status) {
        case Status.SKIP:
          return 'grey lighten-2';
        case Status.PASS:
          return 'green darken-1';
        case Status.WARN:
          return 'warning';
        case Status.FAIL:
          return 'red lighten-1';
        case Status.ERROR:
          return 'error';
        default:
          return 'grey';
      }
    },
    dark() {
      if ([Status.PASS, Status.WARN, Status.ERROR, Status.FAIL].includes(this.status)) {
        return true;
      }

      return false;
    },
  },
});
</script>
