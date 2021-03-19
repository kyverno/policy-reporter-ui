<template>
    <v-autocomplete dense
                    multiple
                    :items="namespaces"
                    outlined
                    hide-details
                    label="Namespaces"
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
import { mapState } from 'vuex';
import debounce from 'lodash.debounce';

const debounced = debounce((emit: () => void) => { emit(); }, 600);

type Data = { selected: string[] }

type Props = { value: string[] }

type Computed = { namespaces: string[] }

type Methods = { input(priorities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
  },
  data: () => ({
    selected: [],
  }),
  computed: mapState(['namespaces']),
  methods: {
    input(namespaces: string[]): void {
      this.selected = namespaces;

      debounced(() => { this.$emit('input', namespaces); });
    },
  },
});
</script>
