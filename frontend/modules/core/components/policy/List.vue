<template>
  <v-toolbar class="my-0" color="secondary">
    <v-toolbar-title class="text-subtitle-1 font-weight-bold">{{ category }}</v-toolbar-title>
    <template #append>
      <Search class="mr-2" v-model="search" style="min-width: 300px;" />
      <CollapseBtn v-model="open" :disabled="!(policies?.length)" />
    </template>
  </v-toolbar>
  <v-list v-if="pending" lines="two" class="mt-0 pt-0 pb-0 mb-0">
    <v-skeleton-loader class="mx-auto border" type="list-item-avatar" />
    <v-skeleton-loader class="mx-auto border" type="list-item-avatar" />
    <v-skeleton-loader class="mx-auto border" type="list-item-avatar" />
  </v-list>
  <template v-else>
    <v-list v-if="list.length && open" lines="two" class="mt-0 pt-0 pb-0 mb-0">
      <policy-list-scroller :list="list" :default-loadings="20" key-prop="name">
        <template #default="{ item }">
          <PolicyItem :item="item" :details="false" :show-status="showed" :summary="summary" />
        </template>
      </policy-list-scroller>
    </v-list>
    <template v-if="!list.length">
      <v-divider />
      <v-card-text>
        No policies for the selected kinds found
      </v-card-text>
    </template>
  </template>
</template>

<script setup lang="ts">
import CollapseBtn from "~/components/CollapseBtn.vue";
import {type PolicyResult, Status} from "~/modules/core/types";

const props = defineProps<{ category: string; policies: PolicyResult[]; pending: boolean; }>()

const search = ref('')
const open = ref(true)

const list = computed(() => {
  if (!search.value) return props.policies

  return props.policies.filter((p) => p.title.toLowerCase().includes(search.value.toLowerCase()))
})


const status = useStatusInjection()

const showed = computed(() => status.value.filter((s) => s !== Status.SKIP))
const summary = computed(() => status.value.includes(Status.SUMMARY))
</script>
