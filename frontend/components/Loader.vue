<template>
  <div v-show="ready">
    <slot />
  </div>

  <v-container fluid v-if="!ready">
    <v-row>
      <v-col class="justify-center align-center text-center d-flex" style="height: 100vh;">
        <v-progress-circular indeterminate size="200" width="10" color="primary" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
const props = defineProps({ time: { type: Number, default: 500 } })

const ready = ref<boolean>(false)

setTimeout(() => { ready.value = true }, props.time)
</script>
