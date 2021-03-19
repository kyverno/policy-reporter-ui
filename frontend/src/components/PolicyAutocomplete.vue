<template>
    <v-autocomplete dense
                    multiple
                    :items="policies"
                    outlined
                    hide-details
                    label="Policies"
                    clearable
                    :value="value"
                    v-bind="$attrs"
                    @input="input"
    >
      <template v-slot:prepend-item>
        <v-list-item ripple @click="toggle">
          <v-list-item-action>
            <v-icon :color="value.length > 0 ? 'indigo darken-4' : ''">
              {{ icon }}
            </v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>
              Select All
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-divider class="mt-2"></v-divider>
      </template>
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

type Props = { policies: string[]; value: string[] }

type Computed = { all: boolean; some: boolean; icon: string }

type Methods = { toggle(): void; input(priorities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    policies: { type: Array, default: () => [] },
    value: { type: Array, default: () => [] },
  },
  data: () => ({
    selected: [],
  }),
  computed: {
    all() {
      return this.value.length === this.policies.length;
    },
    some() {
      return this.value.length > 0 && !this.all;
    },
    icon() {
      if (this.all) return 'mdi-checkbox-marked';
      if (this.some) return 'mdi-minus-box';

      return 'mdi-checkbox-blank-outline';
    },
  },
  methods: {
    input(policies: string[]): void {
      this.selected = policies;

      debounced(() => { this.$emit('input', policies); });
    },
    toggle() {
      this.$nextTick(() => {
        if (this.all) {
          this.$emit('input', []);
        } else {
          const all = this.policies.slice();

          this.selected = all;
          this.$emit('input', all);
        }
      });
    },
  },
  created() {
    if (!this.policies.length || this.value.length > 0) return;

    this.$emit('input', [this.policies[0]]);
  },
});
</script>
