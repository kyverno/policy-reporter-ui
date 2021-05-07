<template>
  <v-dialog width="1024" v-model="open">
    <template v-slot:activator="{ on, attrs }">
      <v-btn dark v-bind="attrs" v-on="on" depressed small> YAML </v-btn>
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

      <v-card-text style="position: relative" v-if="open">
        <highlightjs language="yaml" :code="policy.content" />

        <v-tooltip v-model="show" left>
          <template v-slot:activator="{ attrs }">
            <v-btn
              v-clipboard:copy="policy.content"
              v-clipboard:success="onCopy"
              v-bind="attrs"
              absolute
              style="top: 0px; right: 24px"
              depressed
              class="rounded-0"
              color="grey lighten-2"
            >
              Copy
            </v-btn>
          </template>
          <span>YAML copied</span>
        </v-tooltip>
      </v-card-text>
      <v-card-actions>
          <v-spacer />
          <v-btn text @click="open = false">Close</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import Vue from 'vue';
import { Policy } from '../models';

type Data = { show: boolean; open: boolean };
type Computed = {};
type Props = { policy: Policy };
type Methods = { onCopy(): void };

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    policy: { type: Object, required: true },
  },
  data: () => ({ show: false, open: false }),
  methods: {
    onCopy() {
      this.show = true;

      setTimeout(() => { this.show = false; }, 1500);
    },
  },
});
</script>

<style scoped>
>>> code {
  padding: 16px !important;
}
</style>
