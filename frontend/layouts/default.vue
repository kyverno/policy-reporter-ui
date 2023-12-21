<template>
  <v-layout>
    <v-app-bar elevation="1" prominent>
      <v-app-bar-nav-icon color="white" icon="mdi-menu" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Policy Reporter</v-toolbar-title>
      <template #append>
        <cluster-select />
      </template>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer">
      <v-list density="compact" nav>
        <template v-for="item in navigation" :key="item.title">
          <v-list-item :title="item.title" :to="item.path" :prepend-icon="item.icon" :exact="item.exact"></v-list-item>
        </template>
        <v-divider />
        <template v-for="item in boards" :key="item.title">
          <v-list-item :title="item.title" :to="item.path" :exact="item.exact"></v-list-item>
        </template>
        <v-divider />
        <template v-for="item in sourceNavi" :key="item.path">
          <template v-if="item.children">
            <v-list-item
                :key="item.path"
                :title="item.title"
                :exact="item.exact"
                :to="item.path"
            />
            <v-divider />
            <v-list-item
                v-for="child in item.children"
                :key="child.path"
                :title="child.title"
                :exact="child.exact"
                :to="child.path"
                class="pl-4"
            />
            <v-divider />
          </template>

          <v-list-item :to="item.path" :prepend-icon="item.icon" :exact="item.exact" lines="two" v-else>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
            <v-list-item-subtitle v-if="item.subtitle">{{ item.subtitle }}</v-list-item-subtitle>
          </v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <v-main class="bg-grey-lighten-3">
      <slot />
    </v-main>
  </v-layout>
</template>

<script setup lang="ts">
import ClusterSelect from "~/modules/core/components/select/ClusterSelect.vue";
import type { Source } from "~/modules/core/types";

const drawer = ref(true)

const capitalize = (source: string) => source.charAt(0).toUpperCase() + source.slice(1)

const { data: sources } = useAPI((api) => api.sources(), { default: () => [] })
const { data: customBoards } = useAPI((api) => api.customBoards(), { default: () => [] })

const navigation = [
  { title: 'Dashboard', path: '/', exact: true },
];

const boards = computed(() => (customBoards.value || []).map(b => ({
  title: b.name, path: `/custom-boards/${b.id}`, exact: true
})))

const sourceNavi = computed(() => {
  const list = sources.value as Source[]
  if (list.length === 0) return []

  if (list.length > 1) {
    return list.sort((a, b) => a.name.localeCompare(b.name)).map(s => {
      if (s.categories.length === 0) {
        return { title: capitalize(s.name), path: `/source/${s.name}`, exact: true }
      }

      if (s.categories.length === 1) {
        return { title: capitalize(s.name), subtitle: s.categories[0].name, path: `/source/${s.name}/${s.categories[0].name}`, exact: true }
      }

      return { title: capitalize(s.name), path: `/source/${s.name}`, children: s.categories.map(c => (
            { title: c.name, icon: 'mdi-format-list-checks', path: `/source/${s.name}/${c.name}`, exact: true }
        ))
      }
    })
  }

  return list[0].categories.map(c => ({ title: c.name, icon: 'mdi-format-list-checks', path: `/source/${list[0].name}/${c.name}`, exact: true }))
})
</script>
