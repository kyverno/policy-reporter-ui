<template>
  <v-layout>
    <v-app-bar elevation="1" prominent>
      <v-app-bar-nav-icon icon="mdi-menu" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Policy Reporter</v-toolbar-title>
      <template #append>
        <form-cluster-select />
        <form-display-mode-select style="width: 150px;" />
        <user-menu class="ml-4" :profile="layout.profile" v-if="layout.profile" />
      </template>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer">
      <v-list density="compact" nav color="header" variant="flat">
        <v-list-item title="Dashboard" to="/" exact />
        <v-list-item title="Notification Targets" to="/targets" exact v-if="layout.targets" />

        <v-list-group value="policies" fluid v-if="layout.policies.length > 1">
          <template v-slot:activator="{ props }">
            <v-list-item title="Policies" v-bind="props" base-color="header-item"></v-list-item>
          </template>
          <v-divider class="mb-1" />
          <v-list-item
              :value="child.path"
              base-color="sub-item"
              v-for="child in layout.policies"
              :key="child.path"
              :to="child.path"
              exact
          >
            <v-list-item-title class="pl-2">{{ child.title }}</v-list-item-title>
          </v-list-item>
        </v-list-group>

        <v-list-item title="Policies" :to="layout.policies[0].path" exact v-if="layout.policies.length === 1" />

        <v-divider class="mb-1" />
        <v-list-subheader v-if="layout.customBoards.length">Custom Boards</v-list-subheader>
        <template v-for="item in layout.customBoards" :key="item.title">
          <v-list-item :title="item.title" :to="item.path" :exact="item.exact"></v-list-item>
        </template>
        <v-divider class="mb-1" v-if="layout.customBoards.length" />

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

          <v-list-item :to="item.path" :prepend-icon="item.icon" exact :lines="item.subtitle ? 'two' : 'one'" base-color="header-item" v-else>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
            <v-list-item-subtitle v-if="item.subtitle" style="opacity: 0.75">{{ item.subtitle }}</v-list-item-subtitle>
          </v-list-item>
        </template>
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
import {use} from "h3";
import {useConfigStore} from "~/store/config";

const drawer = ref(true)

const { data: layout } = useAPI((api) => api.layout(), { default: (): LayoutConfig => ({ sources: [], customBoards: [], policies: [] }) })

const theme = useTheme()

const config = useConfigStore()

const bg = computed(() => {
  if (theme.current.value.dark) {
    return 'bg-grey-darken-3'
  }

  return 'bg-grey-lighten-4'
})
</script>
