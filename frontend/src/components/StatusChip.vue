<template>
    <v-chip :color="color" :dark="dark" v-bind="$attrs" v-on="$listeners">{{ status }}</v-chip>
</template>

<script lang="ts">
import { mapStatus } from '@/mapper';
import { Status } from '@/models';
import Vue, { PropType } from 'vue';

export default Vue.extend<{}, {}, { color: string; dark: boolean }, { status: Status }>({
  props: {
    status: { type: String as PropType<Status>, required: true },
  },
  inheritAttrs: false,
  computed: {
    color() {
      return mapStatus(this.status);
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
