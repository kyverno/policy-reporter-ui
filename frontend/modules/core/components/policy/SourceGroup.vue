<template>
  <app-row>
    <v-card>
      <v-toolbar color="category">
        <v-toolbar-title>{{ source.title }}</v-toolbar-title>
        <template #append>
          <CollapseBtn v-model="open"/>
        </template>
      </v-toolbar>
      <div v-show="open">
        <v-divider/>
        <v-card-text>
          <GraphBarPerCategory :source="source.chart" />
        </v-card-text>
        <policy-list v-for="item in source.categories"
                     :key="item"
                     :category="item"
                     :pending="pending as boolean"
                     :policies="data[item] || []"
        />
      </div>
    </v-card>
  </app-row>
</template>

<script lang="ts" setup>
import { APIFilter, NamespacedKinds } from "~/modules/core/provider/dashboard";
import type { Filter, SourceDetails } from "~/modules/core/types";
import CollapseBtn from "~/components/CollapseBtn.vue";
import type { Ref } from "vue";
import { onChange } from "~/helper/compare";

const props = defineProps<{ source: SourceDetails; }>();

const open = ref(true)

const filter = inject<Ref<Filter>>(APIFilter, ref<Filter>({}))
const kinds = inject<Ref<string[]>>(NamespacedKinds, ref<string[]>([]))

const combinedFilter = computed(() => ({
  ...filter.value,
  kinds: kinds.value.length ? kinds.value : undefined,
}))

const { data, refresh, pending } = useAPI(
    (api) => api.policies(props.source.name, combinedFilter.value),
    { default: () => ({}) }
);

watch(combinedFilter, onChange(refresh))

useStatusProvider(ref({ status: props.source.status }))
</script>
