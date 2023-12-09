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
      </v-list>
    </v-navigation-drawer>

    <v-main class="bg-grey-lighten-3">
      <slot />
    </v-main>
  </v-layout>
</template>

<script setup lang="ts">
import ClusterSelect from "~/modules/core/components/select/ClusterSelect.vue";

const drawer = ref(true)

const navigation = [
  { title: 'Dashboard', icon: 'mdi-monitor-dashboard', path: '/', exact: true },
]
</script>
