<template>
  <v-list-group v-bind="$attrs" no-action>
    <template #activator>
      <slot name="icon">
        <v-list-item-icon v-if="icon">
          <v-icon>{{ icon }}</v-icon>
        </v-list-item-icon>
      </slot>

      <v-list-item-content>
        <v-list-item-title>{{ title }}</v-list-item-title>
      </v-list-item-content>
    </template>

    <template v-for="item in items">
      <v-list-item :key="item.title" :to="item.path" :exact="item.exact">
        <v-list-item-icon v-if="item.icon">
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title style="text-transform: capitalize;">
            {{ item.title }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
  </v-list-group>
</template>

<script lang="ts">
import Vue from 'vue'

type Page = { title: string; path: string; exact: boolean; icon?: string }

export default Vue.extend({
  inheritAttrs: false,
  props: {
    icon: { type: String as Vue.PropType<string>, default: '' },
    title: { type: String as Vue.PropType<string>, required: true },
    items: { type: Array as Vue.PropType<Page[]>, required: true }
  }
})
</script>
