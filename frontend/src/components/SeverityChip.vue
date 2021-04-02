<template>
    <v-chip :color="color" dark v-bind="$attrs" v-on="$listeners">{{ severity }}</v-chip>
</template>

<script lang="ts">
import { Severity } from '@/models';
import Vue, { PropType } from 'vue';

export default Vue.extend<{}, {}, { color: string }, { severity: Severity }>({
  props: {
    severity: { type: String as PropType<Severity>, required: true },
  },
  inheritAttrs: false,
  computed: {
    color() {
      switch (this.severity) {
        case Severity.LOW:
          return 'blue lighten-2';
        case Severity.MEDIUM:
          return 'warning';
        case Severity.HIGH:
          return 'red lighten-1';
        default:
          return 'grey';
      }
    },
  },
});
</script>
