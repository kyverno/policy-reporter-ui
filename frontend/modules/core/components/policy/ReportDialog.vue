<template>
    <v-dialog v-model="open" max-width="500">
      <template v-slot:activator="{ props }">
        <v-btn v-bind="props" rounded="4" class="mr-4" variant="tonal" color="white" height="40">Generate Report</v-btn>
      </template>

      <v-card title="Generate HTML Report">
        <v-divider class="mt-2" />
        <v-container>
          <app-row>
            <form-report-select v-model="report" />
          </app-row>
          <app-row v-if="!category">
            <form-category-select :source="source" v-model="categories" />
          </app-row>
          <app-row>
            <form-namespace-select :source="source" v-model="namespaces" />
          </app-row>
          <app-row>
            <policy-kind-select :source="source" v-model="kinds" />
          </app-row>
          <app-row v-if="report === 'policy-report'">
            <v-switch v-model="clusterScope" hide-details label="Cluster Scoped Resources" color="info" />
          </app-row>
        </v-container>
        <v-divider />
        <v-card-actions>
          <v-btn rounded="2" @click="close">Close</v-btn>
          <v-spacer />
          <v-btn rounded="2" @click="request" :loading="loading">Generate</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { callAPI } from "~/modules/core/composables/api";

const props = defineProps<{ source: string; category?: string }>()

const open = ref(false)
const report = ref<string>('policy-report')
const categories = ref<string[]>(props.category ? [props.category] : [])
const namespaces = ref<string[]>([])
const kinds = ref<string[]>([])
const clusterScope = ref<boolean>(true)

const loading = ref<boolean>(false)
const err = ref<Error>()

const close = () => {
  open.value = false

  setTimeout(() => {
    categories.value = []
    namespaces.value = []
  }, 200)
}

const request = async () => {
  loading.value = true
  let response: BlobPart | null = null

  try {
    if (report.value === 'policy-report') {
      response = await callAPI((api) => api.policyHTMLReport(props.source, { categories: categories.value , namespaces: namespaces.value, kinds: kinds.value, clusterScope: clusterScope.value }))
    } else {
      response = await callAPI((api) => api.namespaceHTMLReport(props.source, { categories: categories.value , namespaces: namespaces.value, kinds: kinds.value }))
    }
  } catch (error) {
    err.value = error
    return
  } finally {
    loading.value = false
  }

  const url = window.URL.createObjectURL(new Blob([response], { type: 'text/html; charset=utf-8' }))
  const link = document.createElement('a')

  link.href = url
  link.setAttribute('target', '_blank')

  document.body.appendChild(link)
  link.click()
  URL.revokeObjectURL(url)
}
</script>