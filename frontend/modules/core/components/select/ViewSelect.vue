<template>
  <v-select
    density="compact"
    :items="views"
    variant="outlined"
    hide-details
    label="View"
    :model-value="modelValue"
    @update:model-value="input"
    v-bind="$attrs"
  />
</template>

<script lang="ts" setup>
const props = defineProps({
  modelValue: { type: String, default: "status" },
});

const views = [
  { title: "Group Results by Status", value: "status" },
  { title: "Group Results by Policy", value: "policies" },
  { title: "Group Results by Rule", value: "rules" },
  { title: "Group Results by Category", value: "categories" },
];

const emit = defineEmits(["update:modelValue"]);

const route = useRoute();

if (route.query.views) {
  emit("update:modelValue", route.query.views);
}

const router = useRouter();
const input = (view: string): void => {
  emit("update:modelValue", view);

  router.push({
    name: route.name as string,
    query: { ...route.query, view },
    params: route.params,
  });
};
</script>
