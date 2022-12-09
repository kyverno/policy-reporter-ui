<template>
  <v-select
    dense
    :items="types"
    outlined
    hide-details
    label="Report Type"
    clearable
    :value="selected"
    v-bind="$attrs"
    @input="input"
  />
</template>

<script lang="ts">
import Vue from 'vue'
import debounce from 'lodash.debounce'
import { ReportType } from '../types'

const debounced = debounce((emit: () => void) => { emit() }, 600)

type Data = { types: string[]; selected: ReportType }

type Computed = {}

type Props = { value: ReportType }

type Methods = { input(type: ReportType): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: String as Vue.PropType<ReportType>, default: ReportType.POLICY }
  },
  data: () => ({
    selected: ReportType.POLICY,
    types: [ReportType.POLICY, ReportType.NAMESPACE]
  }),
  watch: {
    value (value: ReportType) {
      if (value !== this.selected) {
        this.input(value)
      }
    }
  },
  created () {
    if (this.$route.query.type) {
      const type = (Array.isArray(this.$route.query.type) ? this.$route.query.type[0] : this.$route.query.type) as ReportType

      this.selected = type
      this.$emit('input', type)
    }
  },
  methods: {
    input (type: ReportType): void {
      this.selected = type

      debounced(() => {
        this.$emit('input', type)
        this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, type } })
      })
    }
  }
})
</script>
