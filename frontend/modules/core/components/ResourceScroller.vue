<template>
  <v-infinite-scroll :onLoad="load" class="no-scrollbar" v-if="loaded.length">
    <template v-for="item in loaded" :key="item">
      <v-row>
        <v-col>
          <slot :item="item" />
        </v-col>
      </v-row>
    </template>
    <template #empty></template>
  </v-infinite-scroll>
</template>

<script setup lang="ts">
import { useInfinite } from "~/composables/infinite";
import { type Filter } from "../types";

const props = defineProps<{ list: any[]; }>()

const list = ref<any[]>(props.list)

watch(() => props.list, () => { list.value = props.list }, { immediate: true })

const { load, loaded } = useInfinite(list)
</script>
