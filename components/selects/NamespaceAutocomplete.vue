<template>
  <v-autocomplete
    dense
    multiple
    :items="namespaces"
    outlined
    hide-details
    label="Namespaces"
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

type Data = { selected: string[]; namespaces: string[]; interval: any }

type Props = { value: string[], source?: string }

type Computed = { refreshInterval: number }

type Methods = { input(priorities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
    source: { type: String, default: undefined }
  },
  data: () => ({
    selected: [],
    namespaces: [],
    interval: null
  }),
  fetch () {
    return this.$coreAPI.namespaces(this.source).then((namespaces) => {
      this.namespaces = namespaces
      this.input(this.selected.filter(s => namespaces.includes(s)))
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
    },
    value (value: string[]) {
      if (value && value.length === 0 && this.selected.length > 0) {
        this.input([])
      }
    }
  },
  destroyed () {
    clearInterval(this.interval)
  },
  created () {
    if (this.$route.query.namespaces) {
      const namespaces = Array.isArray(this.$route.query.namespaces) ? this.$route.query.namespaces.filter(c => !!c) as string[] : [this.$route.query.namespaces]

      this.selected = namespaces
      this.$emit('input', namespaces)
    }
  },
  methods: {
    input (namespaces: string[]): void {
      this.selected = namespaces

      debounced(async () => {
        this.$emit('input', namespaces)
        try {
          await this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, namespaces } })
        } catch {}
      })
    }
  }
})
</script>
