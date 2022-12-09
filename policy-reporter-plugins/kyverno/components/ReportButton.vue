<template>
  <v-btn
    type="button"
    :loading="loading"
    color="info"
    v-bind="$attrs"
    depressed
    @click="request"
  >
    {{ label ? label : (open ? 'Open' : 'Download') }}
  </v-btn>
</template>

<script lang="ts">
import Vue from 'vue'
import { Filter, ReportType } from '../types'

type Data = {
    loading: boolean;
}

type Computed = {
    filter: Filter;
    api: (filter: Filter) => Promise<BlobPart>
}

type Props = {
    policies: string[];
    namespaces: string[],
    clusterScope: boolean;
    type: ReportType;
    open: boolean;
    label: string;
}

type Methods = {
    request: () => Promise<void>
}

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    policies: { type: Array, default: () => [] },
    namespaces: { type: Array, default: () => [] },
    clusterScope: { type: Boolean, default: true },
    open: { type: Boolean, default: false },
    label: { type: String, default: '' },
    type: { type: String as Vue.PropType<ReportType>, required: true }
  },
  data: () => ({ loading: false }),
  computed: {
    filter (): Filter {
      return {
        policies: this.policies,
        namespaces: this.namespaces,
        clusterScope: this.clusterScope ? 1 : 0
      }
    },
    api (): (filter: Filter) => Promise<BlobPart> {
      if (this.type === ReportType.POLICY) {
        return this.$kyvernoAPI.policyReport
      }

      return this.$kyvernoAPI.namespaceReport
    }
  },
  methods: {
    async request () {
      this.loading = true
      let response: BlobPart | null = null

      try {
        if (this.type === ReportType.POLICY) {
          response = await this.$kyvernoAPI.policyReport(this.filter)
        } else {
          response = await this.$kyvernoAPI.namespaceReport(this.filter)
        }
      } catch (error) {
        this.$emit('on-error', error)
        return
      } finally {
        this.loading = false
      }

      const url = window.URL.createObjectURL(new Blob([response], { type: 'text/html; charset=utf-8' }))
      const link = document.createElement('a')
      link.href = url

      if (this.open) {
        link.setAttribute('target', '_blank')
      } else {
        link.setAttribute('download', 'report.html')
      }

      document.body.appendChild(link)
      link.click()
      URL.revokeObjectURL(url)
    }
  }
})
</script>
