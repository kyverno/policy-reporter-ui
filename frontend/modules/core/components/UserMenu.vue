<template>
  <v-menu v-if="profile?.id">
    <template #activator="{ props }">
      <v-btn variant="outlined" class="text-primary ml-3" rounded="4" height="40" v-bind="props">{{ profile.name }}</v-btn>
    </template>
    <v-list class="pb-0">
      <v-list-item @click="logout">
        <v-list-item-title>Logout</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
  <v-chip label variant="outlined" class="text-primary ml-3" rounded="4" size="large" v-if="profile && !profile?.id">{{ profile.name }}</v-chip>
</template>

<script setup lang="ts">
const logout = () => {
  document.cookie.split(";").forEach((c) => { document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/"); });
  // @ts-ignore
  window.location = `${window.location.pathname || '/'}logout`
}

const { data: profile } = useAPI(api => api.profile())

</script>
