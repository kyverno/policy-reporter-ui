<template>
  <div>
    <error-alert v-if="error" class="px-6" />
    <slot v-show="show" />
    <v-container v-if="!show" fluid>
      <v-row>
        <v-col class="justify-center align-center text-center d-flex" style="height: 100vh;">
          <v-progress-circular indeterminate size="200" width="10" />
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend<{ show: boolean }, {}, {}, { loading: boolean; error?: Error }>({
  props: {
    loading: { type: Boolean, default: false },
    error: { type: [Error, Object], default: undefined }
  },
  data: () => ({ show: false }),
  watch: {
    loading: {
      immediate: true,
      handler (loading: boolean) {
        if (!loading) { return }

        setTimeout(() => { this.show = true }, 300)
      }
    }
  }
})
</script>
