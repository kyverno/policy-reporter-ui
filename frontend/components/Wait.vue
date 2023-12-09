<template>
  <div>
    <slot v-if="ready" />
    <slot name="placeholder" v-else />
  </div>
</template>

<script lang="ts" setup>
  const props = defineProps({ time: { type: Number, default: 100 } })

  const ready = ref<boolean>(false)

  setTimeout(() => { ready.value = true }, props.time)
</script>
