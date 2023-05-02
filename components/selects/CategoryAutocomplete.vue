<template>
  <v-autocomplete
    dense
    multiple
    :items="categories"
    outlined
    hide-details
    label="Categories"
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

type Data = { selected: string[]; categories: string[]; interval: any }

type Computed = {}

type Props = { value: string[], source?: string; namespaced: boolean }

type Methods = { input(priorities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    value: { type: Array, default: () => [] },
    source: { type: String, default: undefined },
    namespaced: { type: Boolean, default: false }
  },
  data: () => ({
    interval: null,
    categories: [],
    selected: []
  }),
  fetch () {
    if (this.namespaced) {
      return this.$coreAPI.namespacedCategories(this.source).then((categories) => {
        this.categories = [...categories]

        this.$emit('input', [...(this.selected.length ? categories.filter(s => this.selected.includes(s)) : categories)])
      })
    }

    return this.$coreAPI.clusterCategories(this.source).then((categories) => {
      this.categories = [...categories]

      this.$emit('input', [...(this.selected.length ? categories.filter(s => this.selected.includes(s)) : categories)])
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
    if (this.$route.query.categories) {
      const categories = Array.isArray(this.$route.query.categories) ? this.$route.query.categories.filter(c => !!c) as string[] : [this.$route.query.categories]

      if (categories.length) {
        this.input(categories)
      }
    }
  },
  methods: {
    input (categories: string[]): void {
      this.selected = categories

      debounced(() => {
        this.$emit('input', [...categories])
        this.$router.push({ name: this.$route.name as string, query: { ...this.$route.query, categories } })
      })
    }
  }
})
</script>
