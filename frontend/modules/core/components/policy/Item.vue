<template>
  <v-divider />
  <v-list-item :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }}">
    <template v-slot:prepend>
      <CollapseBtn v-if="item.description" btn-class="mr-2" v-model="open" :size="40" />
      <AvatarSeverity :severity="item.severity ?? Severity.INFO" />
    </template>
    <v-list-item-title>
      {{ item.title }}
    </v-list-item-title>
    <template v-slot:append v-if="summary">
      <ResultChip :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }, query: { kinds: route.query?.kinds, 'cluster-kinds': route.query['cluster-kinds'] }}" class="ml-2" :status="Status.SUMMARY" :count="count" tooltip="results" />
    </template>

    <template v-slot:append v-else>
      <ResultChip v-if="showSkipped" :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }, query: { status: Status.SKIP, kinds: route.query?.kinds, 'cluster-kinds': route.query['cluster-kinds'] }}" class="ml-2" :status="Status.SKIP" :count="item.results[Status.SKIP]" tooltip="skip results" />
      <template v-for="status in showStatus" :key="status">
        <ResultChip :to="{ name: 'policies-source-policy', params: { source: item.source, policy: item.name }, query: { status, kinds: route.query?.kinds, 'cluster-kinds': route.query['cluster-kinds'] }}" class="ml-2" :status="status" :count="item.results[status]" :tooltip="`${status} results`" />
      </template>
    </template>
  </v-list-item>
  <template v-if="open">
    <v-divider />
    <v-list-item :class="`${bg} text-pre-line`">
      {{ item.description }}
    </v-list-item>
  </template>
</template>

<script setup lang="ts">
import { type PropType } from "vue";
import { type PolicyResult, type Filter, Status, Severity } from "~/modules/core/types";

const open = ref(false)

const route = useRoute()

const bg = useBGColor()

const props = defineProps({
  item: { type: Object as PropType<PolicyResult>, required: true },
  details: { type: Boolean, default: false },
  filter: { type: Object as PropType<Filter>, required: false },
  showStatus: { type: Array as PropType<Status[]>, required: true },
  summary: { type: Boolean, default: false },
})

const status = useStatusInjection()

const showSkipped = computed(() => status.value.includes(Status.SKIP) && !!props.item?.results[Status.SKIP])

const count = computed(() => Object.values(props.item?.results || {}).reduce((sum, v) => sum + v, 0))
</script>
