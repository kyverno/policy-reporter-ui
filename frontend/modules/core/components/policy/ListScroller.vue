<template>
  <v-infinite-scroll :onLoad="load" class="no-scrollbar py-0 my-0 policy-scroller" v-if="loaded.length">
    <template v-for="item in loaded" :key="keyProp ? item[keyProp] : item">
      <slot :item="item" />
    </template>
    <template #empty></template>
  </v-infinite-scroll>
</template>

<script setup lang="ts">
const props = defineProps<{ list: any[]; defaultLoadings: number; keyProp?: string }>()

const max = props.list.length

const loaded = ref<any[]>(max ? props.list.slice(0, Math.min(max, props.defaultLoadings)) : [])

const load = ({ done }: any) => {
  const current = loaded.value.length

  loaded.value = [...loaded.value, ...props.list.slice(current, Math.min(max, current + props.defaultLoadings))]

  if (loaded.value.length === props.list.length) {
    done('empty')
  } else {
    done('ok')
  }
}

</script>

<style>
.policy-scroller .v-infinite-scroll__side {
  padding-top: 0;
  margin-top: 0;

  padding-bottom: 0;
  margin-bottom: 0;
}
</style>