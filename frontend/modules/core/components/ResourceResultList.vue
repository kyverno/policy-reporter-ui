<template>
<v-card :title="namespace">
    <v-list v-if="data?.items?.length" lines="two">
      <ResourceResultItem v-for="item in data.items" :key="item.id" :item="item" :details="details" />
    </v-list>
    <template v-if="!pending && !(data.items.length)">
        <v-divider />
        <v-card-text>
            No resources for the selected kinds found
        </v-card-text>
    </template>
</v-card>
</template>

<script setup lang="ts">
import { kinds } from '~/modules/core/store/filter'
import type { CoreAPI } from "~/modules/core/api";

const props = defineProps<{
  namespace: string;
  details: boolean;
}>()

const { $coreAPI } = useNuxtApp()
const pending = ref(true)
const data = ref({ items: [], count: 0 })
const load = async () => {
  try {
    data.value = await ($coreAPI as CoreAPI).namespacedResourceResults({
      namespaces: [props.namespace as string],
      kinds: kinds.value
    })
  } catch (err) {
    console.error(err)
  } finally {
    pending.value = false
  }
}

await load()

watch(kinds, load)
</script>
