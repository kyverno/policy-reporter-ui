<template>
  <v-chip :color="color" dark v-bind="$attrs">
    {{ status }}
  </v-chip>
</template>

<script lang="ts" setup>
import type { PropType } from "vue";
import { Status } from "~/modules/core/types";
import { useTheme } from "vuetify";
import { mapDarkStatus, mapStatus } from "~/modules/core/mapper";

const props = defineProps({
  status: { type: String as PropType<Status>, default: null },
})

const theme = useTheme()

const color = computed(() => theme.current.value.dark ? mapDarkStatus(props.status as Status) : mapStatus(props.status as Status))
</script>
