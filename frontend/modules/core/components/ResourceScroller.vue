<template>
  <v-infinite-scroll :onLoad="load" class="no-scrollbar">
    <template v-for="ns in loaded" :key="ns">
      <v-row>
        <v-col>
          <LazyResourceResultList :namespace="ns" :details="details" :filter="filter" />
        </v-col>
      </v-row>
    </template>
    <template #empty></template>
  </v-infinite-scroll>
</template>

<script setup lang="ts">
import { useInfinite } from "~/composables/infinite";
import { type Filter } from "../types";

const props = defineProps<{ list: any[]; details: boolean; filter?: Filter; }>()

const list = ref<any[]>(props.list)

watch(() => props.list, () => { list.value = props.list }, { immediate: true })

const { load, loaded } = useInfinite(list)
</script>
