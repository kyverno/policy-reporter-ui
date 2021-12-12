<template>
  <v-autocomplete
    dense
    multiple
    :items="severities"
    outlined
    hide-details
    label="Severities"
    clearable
    :value="selected"
    v-bind="$attrs"
    @input="input"
  >
    <template #selection="{ item }">
      <v-chip small label outlined>
        <span>{{ item }}</span>
      </v-chip>
    </template>
  </v-autocomplete>
</template>

<script lang="ts">
import Vue from 'vue'
import debounce from 'lodash.debounce'
import { Severity } from '~/policy-reporter-plugins/core/types'

const debounced = debounce((emit: () => void) => { emit() }, 600)

type Data = { severities: string[]; selected: string[] }

type Computed = {}

type Props = { value: string[] }

type Methods = { input(severities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] }
  },
  data: () => ({
    selected: [],
    severities: [Severity.LOW, Severity.MEDIUM, Severity.HIGH]
  }),
  created () {
    if (this.$route.query.severities) {
      const severities = Array.isArray(this.$route.query.severities) ? this.$route.query.severities.filter(c => !!c) as string[] : [this.$route.query.severities]

      this.selected = severities
      this.$emit('input', severities)
    }
  },
  methods: {
    input (severities: string[]): void {
      this.selected = severities

      debounced(() => {
        this.$emit('input', severities)
        this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, severities } })
      })
    }
  }
})
</script>
