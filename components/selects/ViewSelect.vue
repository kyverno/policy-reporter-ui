<template>
  <v-select
    dense
    :items="views"
    outlined
    hide-details
    label="View"
    :value="value"
    v-bind="$attrs"
    @input="input"
  />
</template>

<script lang="ts">
import Vue from 'vue'

type Data = { views: Array<{ text: string; value: string }> }

type Computed = {}

type Props = { value: string }

type Methods = { input(view: string): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: String, default: 'status' }
  },
  data: () => ({
    views: [
      { text: 'Group Results by Status', value: 'status' },
      { text: 'Group Results by Policy', value: 'policies' },
      { text: 'Group Results by Category', value: 'categories' }
    ]
  }),
  created () {
    if (this.$route.query.view) {
      this.$emit('input', this.$route.query.view)
    }
  },
  methods: {
    input (view: string): void {
      this.$emit('input', view)
      this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, view } })
    }
  }
})
</script>
