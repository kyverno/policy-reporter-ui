<template>
  <v-infinite-scroll :onLoad="load" class="no-scrollbar" v-if="loaded.length">
    <template v-for="item in loaded" :key="item">
      <slot :item="item" />
    </template>
    <template #empty></template>
  </v-infinite-scroll>
</template>

<script setup lang="ts">
import { useInfinite } from "~/composables/infinite";

const props = defineProps<{ list: any[]; }>()

const list = ref<any[]>(props.list)

watch(() => props.list, () => { list.value = props.list }, { immediate: true })

const { load, loaded } = useInfinite(list, props.list.length)
</script>
