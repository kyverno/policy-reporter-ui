<template>
  <v-btn :color="color" clss="py-1" size="small" :value="status">{{ status }}</v-btn>
</template>

<script setup lang="ts">
import { Status } from "~/modules/core/types";
import { useTheme } from "vuetify";
import { useStatusBtnColor } from "~/modules/core/composables/theme";

const props = defineProps<{ status: Status }>()

const theme = useTheme()

const color = useStatusBtnColor(props.status)

</script>
