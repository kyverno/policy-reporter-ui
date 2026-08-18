<template>
  <v-data-table-server
    :items="results.results"
    :items-length="results.count"
    :headers="headers"
    item-value="id"
    show-expand
    :items-per-page="itemsPerPage"
    :page="page"
    @update:items-per-page="emit('update:itemsPerPage', $event)"
    @update:page="emit('update:page', $event)"
    hover
  >
    <template #item="{ item, ...props }">
      <tr @click="() => props.toggleExpand(props.internalItem)" class="cursor-pointer">
        <td>
          <v-btn v-if="plugin" @click.stop="openPolicy(item.policy)" class="mr-1" target="_blank" icon="mdi-open-in-new" variant="text" size="small" />
          {{ item.policy }}
        </td>
        <td>{{ item.rule }}</td>
        <td>
          <chip-severity v-if="item.severity" @click.prevent.stop="emit('search', item.severity)" :severity="item.severity" />
        </td>
        <td>
          <chip-status @click.prevent.stop="emit('search', item.status)" :status="item.status" />
        </td>
        <td>
          <resource-exception-dialog v-if="exceptions" :resource="resource" :source="source" :policies="[{ name: item.policy, rules: [{ name: item.rule, props: item.properties }]}]" />
        </td>
      </tr>
    </template>
    <template #expanded-row="{ columns, item }">
      <tr :class="bg">
        <td :colspan="columns.length" class="py-3">
          <div v-if="item.hasProps">
            <v-card v-if="item.message" variant="flat">
              <v-alert type="info" variant="flat" class="text-pre-line">
                {{ item.message }}
              </v-alert>
            </v-card>
            <div>
              <template v-for="(value, label) in item.chips"  :key="label">
                <property-chip :label="label as string" :value="value" class="mr-2 mt-2 rounded-lg" />
              </template>
              <template v-for="(value, label) in item.cards"  :key="label">
                <property-card :label="label as string" :value="value" class="mt-2" />
              </template>
            </div>
          </div>
          <div class="text-pre-line" v-else>
            {{ item.message }}
          </div>
        </td>
      </tr>
    </template>
  </v-data-table-server>
</template>

<script lang="ts" setup>
import { type ResultList } from "~/types/core";

import { mapResults } from "~/core/mapper";

const props = defineProps<{
  plugin?: boolean;
  source: string;
  exceptions?: boolean;
  resource: string;
  data: ResultList;
  itemsPerPage?: number;
  page?: number;
}>()

const emit = defineEmits<{
  (e: 'search', value: string): void
  (e: 'update:itemsPerPage', value: number): void
  (e: 'update:page', value: number): void
}>()

const bg = useBGColor()

const router = useRouter()
const openPolicy = (policy: string) => {
  const { href } = router.resolve({ name: 'policies-source-info-policy', params: { source: props.source, policy }})
  window.open(href, '_blank')
}

const results = computed(() => mapResults(props.data as ResultList))

const headers = computed(() => {
  if (props.exceptions) {
    return [
      { title: 'Policy', key: 'policy', width: '33%' },
      { title: 'Rule', key: 'rule', width: '33%' },
      { title: 'Severity', key: 'severity', width: '12%' },
      { title: 'Status', key: 'status', width: '12%' },
      { title: 'Actions', key: 'exception', width: '12%', sortable: false }
    ]
  }

  return [
    { title: 'Policy', key: 'policy', width: '33%' },
    { title: 'Rule', key: 'rule', width: '33%' },
    { title: 'Severity', key: 'severity', width: '17%' },
    { title: 'Status', key: 'status', width: '17%' }
  ]
})
</script>
