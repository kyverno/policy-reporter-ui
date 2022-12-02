<template>
  <v-autocomplete
    dense
    multiple
    :items="labels"
    outlined
    hide-details
    :label="capitalize(name)"
    clearable
    :value="selected"
    v-bind="$attrs"
    @input="input"
  >
    <template #selection="{ item, index }">
      <v-chip v-if="index <= 1" small label outlined>
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
import Vue from 'vue'
import debounce from 'lodash.debounce'

const debounced = debounce((emit: () => void) => { emit() }, 600)

type Data = { selected: string[]; interval: any }

type Computed = {}

type Props = { value: string[], source?: string; labels: string[]; name: string }

type Methods = {
  input(priorities: string[]): void;
  capitalize (input: string): string;
}

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
    labels: { type: Array, default: () => [] },
    name: { type: String, default: undefined }
  },
  data: () => ({
    interval: null,
    selected: []
  }),
  created () {
    if (this.$route.query.labels) {
      const labels = (Array.isArray(this.$route.query.labels) ? (this.$route.query.labels as string[]).filter(c => !!c) as string[] : [this.$route.query.labels as string])
        .filter(l => l.startsWith(`${this.name}:`))
        .map(l => l.replace(`${this.name}:`, ''))

      if (labels.length) {
        this.input(labels)
      }
    }
  },
  methods: {
    input (labels: string[]): void {
      this.selected = labels

      debounced(() => {
        this.$emit('input', [...labels])
        this.$router.push({
          name: this.$route.name as string,
          query: {
            ...this.$route.query,
            labels: [
              ...(Array.isArray(this.$route.query.labels) ? this.$route.query.labels as string[] : [this.$route.query.labels] as string[]).filter(l => l && !l.startsWith(`${this.name}:`)),
              ...labels.map(l => `${this.name}:${l}`)
            ]
          }
        })
      })
    },
    capitalize (input: string): string {
      return input[0].toUpperCase() + input.slice(1)
    }
  }
})
</script>
