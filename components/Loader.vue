<template>
  <div>
    <error-alert v-if="error" class="px-6" />
    <div v-show="show" :key="key">
      <slot />
    </div>
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
import { mapGetters } from 'vuex'
import { Cluster } from '~/policy-reporter-plugins/core/types'

export default Vue.extend<{ show: boolean }, { currentCluster?: Cluster; key: string }, {}, { loading: boolean; error?: Error }>({
  props: {
    loading: { type: Boolean, default: false },
    error: { type: [Error, Object], default: undefined }
  },
  data: () => ({ show: false }),
  computed:
  {
    ...mapGetters(['currentCluster']),
    key (): string {
      if (this.currentCluster) {
        return this.$route.name + this.currentCluster.id
      }

      return this.$route.name as string
    }
  },
  watch: {
    loading: {
      immediate: true,
      handler (loading: boolean) {
        if (loading) {
          this.show = false
          return
        }

        setTimeout(() => {
          this.$nextTick(() => { this.show = true })
        }, 300)
      }
    }
  }
})
</script>
