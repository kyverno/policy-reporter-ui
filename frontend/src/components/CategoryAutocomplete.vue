<template>
    <v-autocomplete dense
                    multiple
                    :items="categories"
                    outlined
                    hide-details
                    label="Categories"
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

type Props = { categories: string[]; value: string[] }

type Methods = { input(priorities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
    categories: { type: Array, default: () => [] },
  },
  data: () => ({
    selected: [],
  }),
  watch: {
    categories(categories: string[], before: string[]) {
      if (JSON.stringify(categories.sort()) === JSON.stringify(before.sort())) return;

      this.$emit('input', this.value.filter((value) => categories.includes(value)));
    },
  },
  methods: {
    input(categories: string[]): void {
      this.selected = categories;

      debounced(() => {
        this.$emit('input', categories);
        this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, categories } });
      });
    },
  },
  created() {
    if (this.$route.query.categories) {
      const categories = Array.isArray(this.$route.query.categories) ? this.$route.query.categories.filter((c) => !!c) as string[] : [this.$route.query.categories];

      this.selected = categories;
      this.$emit('input', categories);
    }
  },
});
</script>
