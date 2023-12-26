<template>
  <wait>
      <v-toolbar class="my-0" color="secondary">
        <v-toolbar-title class="text-subtitle-1 font-weight-bold">{{ category }}</v-toolbar-title>
        <template #append>
          <Search class="mr-2" v-model="search" style="min-width: 300px;" />
          <CollapseBtn v-model="open" :disabled="!data.length" />
        </template>
      </v-toolbar>
      <v-list v-if="data?.length && open" lines="two" class="mt-0 pt-0">
        <PolicyItem v-for="item in data" :key="item.policy" :item="item" :details="false" />
      </v-list>
      <template v-if="!pending && !(data?.length)">
        <v-divider />
        <v-card-text>
          No policies for the selected kinds found
        </v-card-text>
      </template>
  </wait>
</template>

<script setup lang="ts">
import type { Ref } from "vue";
import CollapseBtn from "~/components/CollapseBtn.vue";
import type { Filter } from "~/modules/core/types";
import { NamespacedKinds, ResourceFilter } from "~/modules/core/provider/dashboard";
import { execOnChange } from "~/helper/compare";

const props = defineProps<{ category: string; }>()

const search = ref('')
const open = ref(true)

const filter = inject<Ref<Filter>>(ResourceFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(NamespacedKinds, ref<string[]>([]))

const combinedFilter = computed(() => ({
  ...filter.value,
  categories: [props.category as string],
  kinds: kinds.value.length ? kinds.value : undefined,
  search: search.value,
}))

const { data, refresh, pending } = useAPI(
    (api) => api.policies(combinedFilter.value),
    { default: () => ([]) }
);

watch(combinedFilter, (o, n) => execOnChange(o, n, () => refresh()))
</script>
