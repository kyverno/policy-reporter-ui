<template>
    <v-dialog v-model="open" max-width="700">
      <template v-slot:activator="{ props }">
        <v-btn @click.prevent.stop="request" rounded="0" class="ml-4" variant="flat" color="secondary" size="small" :height="height as any">Exception</v-btn>
      </template>

      <v-card title="Resource Exception">
        <v-divider class="mt-2" />
        <v-container>
          <app-row v-if="err">
            <v-alert variant="tonal" type="error">Failed to create exception: {{ err }}</v-alert>
          </app-row>
          <template v-else>
            <app-row>
              <v-alert type="warning" variant="tonal">
                Creating many small PolicyExceptions can impact the performance. If you need to exclude multiple resources from a policy consider to extend a single exception.
              </v-alert>
            </app-row>
            <app-row style="position: relative;">
              <highlightjs :code="content" />
              <v-btn alt="Copy to Clipboard" theme="dark" style="position: absolute; top: 25px; right: 25px;" rounded="0" color="primary" variant="tonal" @click="copy(content)" :icon="copied ? 'mdi-content-save-check' : 'mdi-content-save'" />
              <v-btn alt="Download as File" theme="dark" style="position: absolute; top: 25px; right: 87px;" rounded="0" color="primary" variant="tonal" @click="download()" icon="mdi-file-download" />
            </app-row>
          </template>
        </v-container>
        <v-divider />
        <v-card-actions>
          <v-btn rounded="2" @click="close">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { callAPI } from "~/modules/core/composables/api";
import { useClipboard } from '@vueuse/core'
import { FetchError } from "ofetch";
import { parse } from "yaml";


const props = defineProps<{
  source: string;
  resource: string;
  policies?: { name: string; rules: string[] }[];
  height?: string | number;
}>()

const content = ref('')
const open = ref(false)
const loading = ref(false)
const err = ref<string>()

const { text, copy, copied } = useClipboard({ source: content })

const close = () => {
  open.value = false
}

const request = async () => {
  loading.value = true

  try {
    const response = await callAPI((api) => api.createException(props.resource, props.source, props.policies))
    content.value = response.resource
    err.value = undefined

  } catch (error: FetchError) {
    err.value = `[${error.statusCode}] ${error.statusMessage}`
    return
  } finally {
    loading.value = false
    open.value = true
  }
}

const download = () => {
  if (!content.value) return;

  try {
    const res = parse(content.value)

    let element = document.createElement('a');
    element.setAttribute('href', 'data:application/yaml;charset=utf-8,' + encodeURIComponent(content.value));
    element.setAttribute('download', `${res?.metadata?.name || props.resource}.yaml`);

    element.style.display = 'none';
    document.body.appendChild(element);

    element.click();
    document.body.removeChild(element);
  } catch {
  }
}
</script>