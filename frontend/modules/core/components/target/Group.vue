
<template>
  <v-card>
    <v-toolbar color="category">
      <v-toolbar-title>{{ capilize(type) }}</v-toolbar-title>
      <template #append>
        <CollapseBtn v-model="open" />
      </template>
    </v-toolbar>
    <template v-if="open" v-for="target in targets" :key="target.name">
      <v-toolbar color="secondary">
        <v-toolbar-title>{{ target.name }}</v-toolbar-title>
        <template #append>
          <chip-priority :priority="target.minimumPriority" variant="flat" class="mr-2" />
        </template>
      </v-toolbar>
      <v-divider />
      <v-list class="mt-0 pt-0" lines="two">
        <v-list-item v-if="target.host" :title="target.host" subtitle="Host" />
        <template  v-if="target.useTLS">
          <v-divider />
          <v-list-item title="TLS configured" />
        </template>
        <template  v-if="target.auth">
          <v-divider />
          <v-list-item title="Uses Authentication" subtitle="Auth method like HTTP Basic, Authorized Header or API Key configured." />
        </template>
        <template  v-if="target.skipTLS">
          <v-divider />
          <v-list-item title="Skips TLS verification" subtitle="SkipTLS enabled" />
        </template>
        <template  v-for="(v, k) in target.properties">
          <v-divider />
          <v-list-item :title="v" :subtitle="capilize(k)" />
        </template>
        <template  v-if="target.secretRef">
          <v-divider />
          <v-list-item :title="`SecretRef: ${target.secretRef}`" subtitle="Secret used to retrieve sensitive target configuration" />
        </template>
        <template  v-if="target.mountedSecret">
          <v-divider />
          <v-list-item :title="`MountedSecret: ${target.mountedSecret}`" subtitle="Secret mount used to retrieve sensitive target configuration" />
        </template>
      </v-list>
      <template v-if="target.filter && Object.keys(target.filter).length > 0">
        <v-divider />
        <v-table hover>
          <thead>
          <tr>
            <th>Filter</th>
            <th>Included Values</th>
            <th>Excluded Values</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="(values, filter) in target.filter">
            <td>{{ filter }}</td>
            <td>{{ values?.include?.join(', ') }}</td>
            <td>{{ values?.exclude?.join(', ') }}</td>
          </tr>
          </tbody>
        </v-table>
      </template>
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { capilize } from "~/modules/core/layouthHelper";
import { type Target } from "~/modules/core/types";
import CollapseBtn from "~/components/CollapseBtn.vue";

defineProps<{ type: string; targets: Target[]; }>()

const open = ref(true)
</script>
