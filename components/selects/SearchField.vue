<template>
  <v-text-field
    :value="search"
    append-icon="mdi-magnify"
    label="Search"
    outlined
    dense
    hide-details
    @input="input"
  />
</template>

<script lang="ts">
import Vue from 'vue'
import debounce from 'lodash.debounce'

const debounced = debounce((emit: () => void) => { emit() }, 600)

type Data = { search: string; }

type Computed = {}

type Props = { value: string }

type Methods = { input(search: string): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: String, default: '' }
  },
  data: () => ({ search: '' }),
  watch: {
    value: {
      immediate: true,
      handler (search: string) {
        this.search = search
      }
    }
  },
  methods: {
    input (search: string): void {
      this.search = search

      debounced(() => {
        this.$emit('input', search)
      })
    }
  }
})
</script>
