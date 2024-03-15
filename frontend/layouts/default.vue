<template>
  <v-layout v-if="layout">
    <v-app-bar elevation="1" prominent>
      <v-app-bar-nav-icon icon="mdi-menu" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Policy Reporter</v-toolbar-title>
      <template #append>
        <form-cluster-select />
        <form-display-mode-select />
        <user-menu />
      </template>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer">
      <v-list density="compact" nav color="header" variant="flat">
        <v-list-item title="Dashboard" to="/" base-color="header" exact />

        <template v-for="item in layout.sources" :key="item.path">
          <template v-if="item.children">
            <v-list-item
                base-color="header-item"
                :key="item.path"
                :title="item.title"
                :to="item.path"
                exact
            />
            <v-divider class="mb-1" />
            <v-list-item
                base-color="sub-item"
                v-for="child in item.children"
                :key="child.path"
                :title="child.title"
                :to="child.path"
                class="pl-4"
                exact
            />
            <v-divider class="mb-1" />
          </template>

          <v-list-item :to="item.path" exact lines="one" base-color="header-item" v-else>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
        </template>

        <div class="mb-1 mt-8" />

          <v-list-item base-color="header" title="Policy Dashboard" to="/policies" exact />
          <template v-for="item in layout.policies" :key="item.path">
            <v-list-item :to="item.path" exact lines="one" base-color="header-item">
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item>
          </template>

        <div class="mb-1 mt-8" />

        <v-list-item base-color="header" v-if="layout.customBoards.length" title="Custom Boards" />
        <template v-for="item in layout.customBoards" :key="item.path">
          <v-list-item base-color="header-item" :title="item.title" :to="item.path" :exact="item.exact"></v-list-item>
        </template>

        <div class="mb-1 mt-8" />

        <v-list-item title="Notification Targets" to="/targets" exact base-color="header-item" v-if="layout.targets" />
      </v-list>
    </v-navigation-drawer>

    <v-main :class="bg">
      <v-container fluid v-if="config.error">
        <app-row>
          <v-card class="pa-2">
            <v-alert variant="outlined" type="error">Failed to access API: {{ config.error }}</v-alert>
          </v-card>
        </app-row>
      </v-container>

      <slot v-else />
    </v-main>
  </v-layout>
</template>

<script setup lang="ts">
import { useTheme } from "vuetify";
import type { LayoutConfig } from "~/modules/core/types";
import { useConfigStore } from "~/store/config";

const drawer = ref(true)

const { data: layout } = useAPI((api) => api.layout(), { default: (): LayoutConfig => ({ sources: [], customBoards: [], policies: [], targets: false }) })

const theme = useTheme()

const config = useConfigStore()

const bg = computed(() => {
  if (theme.current.value.dark) {
    return 'bg-grey-darken-3'
  }

  return 'bg-grey-lighten-4'
})
</script>
