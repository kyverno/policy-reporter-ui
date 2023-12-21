<template>
  <wait>
    <v-card>
      <v-toolbar color="transparent">
        <v-toolbar-title>{{ namespace }}</v-toolbar-title>
        <template #append>
          <Search class="mr-2" v-model="search" style="min-width: 300px;" />
          <CollapseBtn v-model="open" :disabled="!data.items.length" />
        </template>
      </v-toolbar>
      <v-list v-if="data?.items?.length && open" lines="two">
        <ResourceResultItem v-for="item in data.items" :key="item.id" :item="item" :details="details" :filter="filter" />
      </v-list>
      <template v-if="!pending && !(data.items.length)">
        <v-divider />
        <v-card-text>
          No resources for the selected kinds found
        </v-card-text>
      </template>
    </v-card>
  </wait>
</template>

<script setup lang="ts">
import { kinds } from '~/modules/core/store/filter'
import type { CoreAPI } from "~/modules/core/api";
import CollapseBtn from "~/components/CollapseBtn.vue";
import type { Filter } from "~/modules/core/types";

const props = defineProps<{
  namespace: string;
  details: boolean;
  filter?: Filter;
}>()

const search = ref('')
const open = ref(true)

const { $coreAPI } = useNuxtApp()
const pending = ref(true)
const data = ref({ items: [], count: 0 })
const load = async () => {
  try {
    data.value = await ($coreAPI as CoreAPI).namespacedResourceResults({
      ...(props.filter || {}),
      namespaces: [props.namespace as string],
      kinds: kinds.value,
      search: search.value,
    })
  } catch (err) {
    console.error(err)
  } finally {
    pending.value = false
  }
}

load()

watch(() => ({ kinds: kinds.value, search: search.value }), load)
</script>
