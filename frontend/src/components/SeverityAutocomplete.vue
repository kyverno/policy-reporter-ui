<template>
    <v-autocomplete dense
                    multiple
                    :items="severities"
                    outlined
                    hide-details
                    label="Severities"
                    clearable
                    :value="value"
                    v-bind="$attrs"
                    @input="input"
    >
    <template v-slot:selection="{ item }">
        <v-chip small label outlined>
          <span>{{ item }}</span>
        </v-chip>
      </template>
    </v-autocomplete>
</template>

<script lang="ts">
import Vue from 'vue';
import debounce from 'lodash.debounce';
import { Severity } from '@/models';

const debounced = debounce((emit: () => void) => { emit(); }, 600);

type Data = { severities: string[]; selected: string[] }

type Computed = {}

type Props = { value: string[] }

type Methods = { input(severities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
  },
  data: () => ({
    selected: [],
    severities: [Severity.LOW, Severity.MEDIUM, Severity.HIGH],
  }),
  methods: {
    input(severities: string[]): void {
      this.selected = severities;

      debounced(() => { this.$emit('input', severities); });
    },
  },
});
</script>
