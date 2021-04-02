<template>
    <v-chip :color="color" :dark="dark" v-on="$listeners" v-bind="$attrs">{{ priority }}</v-chip>
</template>

<script lang="ts">
import { Priority } from '@/models';
import Vue, { PropType } from 'vue';

export default Vue.extend<{}, {}, { color: string; dark: boolean }, { priority: Priority }>({
  props: {
    priority: { type: String as PropType<Priority>, required: true },
  },
  computed: {
    color() {
      switch (this.priority) {
        case Priority.DEBUG:
          return 'grey lighten-2';
        case Priority.INFO:
          return 'light-blue lighten-1';
        case Priority.WARNING:
          return 'warning';
        case Priority.ERROR:
          return 'error';
        case Priority.CRITICAL:
          return 'red darken-3';
        default:
          return 'grey';
      }
    },
    dark() {
      if ([Priority.WARNING, Priority.ERROR].includes(this.priority)) {
        return true;
      }

      return false;
    },
  },
});
</script>
