<template>
  <v-divider />
  <v-list-item :to="{ name: 'resource-id', query: filter, params: { id: item.id }}">
    <template v-if="details" v-slot:prepend>
      <v-btn class="mr-2" variant="text" :icon="!open ? `mdi-chevron-up` : `mdi-chevron-down`" @click.stop.prevent="open = !open"></v-btn>
    </template>
    <v-list-item-title>
      {{ item.name }}
<!--      <nuxt-link to="/" target="_blank" class="text-decoration-none" title="open details">-->
<!--        {{ item.name }}-->
<!--      </nuxt-link>-->
    </v-list-item-title>
    <v-list-item-subtitle>{{ item.apiVersion }} {{ item.kind }}</v-list-item-subtitle>
    <template v-slot:append>
      <ResultChip :status="Status.PASS" :count="item.pass" tooltip="pass results" />
      <ResultChip class="ml-2" :status="Status.WARN" :count="item.warn" tooltip="warning results" />
      <ResultChip class="ml-2" :status="Status.FAIL" :count="item.fail" tooltip="fail results" />
      <ResultChip class="ml-2" :status="Status.ERROR" :count="item.error" tooltip="error results" />
    </template>
  </v-list-item>
  <LazyResourceResults v-if="open" :id="item.id" :filter="filter" />
</template>

<script setup lang="ts">
import { type Filter, type ResourceResult, Status } from '../types'
import { type PropType } from "vue";

const open = ref(false)

const props = defineProps({
  item: { type: Object as PropType<ResourceResult>, required: true },
  details: { type: Boolean, default: false },
  filter: { type: Object as PropType<Filter>, required: false },
})

</script>
