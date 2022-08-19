<template>
  <v-autocomplete
    dense
    multiple
    :items="policies"
    outlined
    hide-details
    label="Policies"
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
import { mapGetters } from 'vuex'
import debounce from 'lodash.debounce'

const debounced = debounce((emit: () => void) => { emit() }, 600)

type Data = { selected: string[]; policies: string[]; interval: any }

type Props = { value: string[]; namespaced: boolean; source?: string }

type Computed = {}

type Methods = { input(priorities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
    namespaced: { type: Boolean, default: false },
    source: { type: String, default: undefined }
  },
  data: () => ({
    selected: [],
    policies: [],
    interval: null
  }),
  fetch () {
    if (this.namespaced) {
      return this.$coreAPI.namespacedPolicies(this.source).then((policies) => {
        this.policies = policies

        this.$emit('input', [...(this.selected.length ? policies.filter(s => this.selected.includes(s)) : policies)])
      })
    }

    return this.$coreAPI.clusterPolicies(this.source).then((policies) => {
      this.policies = policies

      this.$emit('input', [...(this.selected.length ? policies.filter(s => this.selected.includes(s)) : policies)])
    })
  },
  computed: mapGetters(['refreshInterval']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  created () {
    if (this.$route.query.policies) {
      const policies = Array.isArray(this.$route.query.policies) ? this.$route.query.policies.filter(c => !!c) as string[] : [this.$route.query.policies]

      if (policies.length) {
        this.input(policies)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  },
  methods: {
    input (policies: string[]): void {
      this.selected = policies

      debounced(() => {
        this.$emit('input', [...policies])
        this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, policies, all: undefined } })
      })
    }
  }
})
</script>
