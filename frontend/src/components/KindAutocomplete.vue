<template>
    <v-autocomplete dense
                    multiple
                    :items="kinds"
                    outlined
                    hide-details
                    label="Kinds"
                    clearable
                    :value="value"
                    v-bind="$attrs"
                    @input="input"
    >
    <template v-slot:selection="{ item, index }">
        <v-chip small v-if="index <= 1" label outlined>
          <span>{{ item }}</span>
        </v-chip>
        <span
          v-if="index === 2"
          class="grey--text caption"
        >
          (+{{ selected.length - 2 }} others)
        </span>
      </template>
    </v-autocomplete>
</template>

<script lang="ts">
import Vue from 'vue';
import debounce from 'lodash.debounce';

const debounced = debounce((emit: () => void) => { emit(); }, 600);

type Data = { selected: string[] }

type Computed = {}

type Props = { kinds: string[]; value: string[] }

type Methods = { input(kinds: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
    kinds: { type: Array, default: () => [] },
  },
  data: () => ({
    selected: [],
  }),
  watch: {
    categories(kinds: string[], before: string[]) {
      if (JSON.stringify(kinds.sort()) === JSON.stringify(before.sort())) return;

      this.$emit('input', this.value.filter((value) => kinds.includes(value)));
    },
  },
  methods: {
    input(kinds: string[]): void {
      this.selected = kinds;

      debounced(() => { this.$emit('input', kinds); });
    },
  },
});
</script>
