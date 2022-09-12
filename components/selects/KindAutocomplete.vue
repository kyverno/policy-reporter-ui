<template>
  <v-autocomplete
    dense
    multiple
    :items="kinds"
    outlined
    hide-details
    label="Kinds"
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

type Data = { selected: string[]; kinds: string[]; interval: any }

type Computed = {}

type Props = { value: string[]; namespaced: boolean, source?: string }

type Methods = { input(kinds: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
    namespaced: { type: Boolean, default: false },
    source: { type: String, default: undefined }
  },
  data: () => ({
    selected: [],
    kinds: [],
    interval: null
  }),
  fetch () {
    if (this.namespaced) {
      return this.$coreAPI.namespacedKinds(this.source).then((kinds) => {
        this.kinds = kinds
        this.input(this.selected.filter(s => kinds.includes(s)))
      })
    }

    return this.$coreAPI.clusterKinds(this.source).then((kinds) => {
      this.kinds = kinds
      this.input(this.selected.filter(s => kinds.includes(s)))
    })
  },
  computed: mapGetters(['refreshInterval']),
  watch: {
    refreshInterval: {
      immediate: true,
      handler (refreshInterval: number) {
        if (this.interval) { clearInterval(this.interval) }

        if (!refreshInterval) {
          this.interval = null
          return
        }

        this.interval = setInterval(this.$fetch, refreshInterval)
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  },
  created () {
    if (this.$route.query.kinds) {
      const kinds = Array.isArray(this.$route.query.kinds) ? this.$route.query.kinds.filter(c => !!c) as string[] : [this.$route.query.kinds]

      this.selected = kinds
      this.$emit('input', kinds)
    }
  },
  methods: {
    input (kinds: string[]): void {
      this.selected = kinds

      debounced(() => {
        this.$emit('input', kinds)
        this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, kinds } })
      })
    }
  }
})
</script>
