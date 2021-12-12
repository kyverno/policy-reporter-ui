<template>
  <v-dialog v-model="open" width="1024">
    <template #activator="{ on, attrs }">
      <v-btn
        dark
        class="action-button"
        v-bind="attrs"
        depressed
        small
        v-on="on"
      >
        YAML
      </v-btn>
    </template>

    <v-card min-height="500">
      <v-card-title>
        {{ policy.name }}
        <v-spacer />
        <v-btn icon @click="open = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-subtitle class="mt-1">
        {{ policy.description }}
      </v-card-subtitle>

      <v-card-text v-if="open" style="position: relative;">
        <yaml-content>
          {{ policy.content }}
        </yaml-content>

        <v-tooltip v-model="show" left>
          <template #activator="{ attrs }">
            <v-btn
              v-clipboard:copy="policy.content"
              v-clipboard:success="onCopy"
              v-bind="attrs"
              absolute
              style="top: 0; right: 24px;"
              depressed
              class="rounded-0 grey-background"
            >
              Copy
            </v-btn>
          </template>
          <span>YAML copied</span>
        </v-tooltip>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="open = false">
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import Vue from 'vue'
import { Policy } from '../types'
import YamlContent from './YamlContent.vue'

type Data = { show: boolean; open: boolean };
type Computed = {};
type Props = { policy: Policy };
type Methods = { onCopy(): void };

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { YamlContent },
  props: {
    policy: { type: Object, required: true }
  },
  data: () => ({ show: false, open: false }),
  methods: {
    onCopy () {
      this.show = true

      setTimeout(() => { this.show = false }, 1500)
    }
  }
})
</script>

<style scoped>
>>> code {
  padding: 16px !important;
}
</style>
