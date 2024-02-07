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
        <template v-for="item in navigation" :key="item.title">
          <v-list-item :title="item.title" :to="item.path" :exact="item.exact" />
        </template>

        <v-list-group value="policies" fluid>
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

        <v-divider class="mb-1" />
        <template v-if="layout.customBoards">
          <v-list-subheader>Custom Boards</v-list-subheader>
        </template>
        <template v-for="item in layout.customBoards" :key="item.title">
          <v-list-item :title="item.title" :to="item.path" :exact="item.exact"></v-list-item>
        </template>
        <v-divider class="mb-1" />
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

          <v-list-item :to="item.path" :prepend-icon="item.icon" exact lines="two" base-color="header-item" v-else>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
            <v-list-item-subtitle v-if="item.subtitle" style="opacity: 0.75">{{ item.subtitle }}</v-list-item-subtitle>
          </v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <v-main :class="bg">
      <slot />
    </v-main>
  </v-layout>
</template>

<script setup lang="ts">
import { useTheme } from "vuetify";
import type { LayoutConfig } from "~/modules/core/types";

const drawer = ref(true)

const { data: layout } = useAPI((api) => api.layout(), { default: (): LayoutConfig => ({ sources: [], customBoards: [], policies: [] }) })

const theme = useTheme()

const bg = computed(() => {
  if (theme.current.value.dark) {
    return 'bg-grey-darken-3'
  }

  return 'bg-grey-lighten-4'
})

const navigation = [
  { title: 'Dashboard', path: '/', exact: true },
  { title: 'Notification Targets', path: '/targets', exact: true },
];
</script>
