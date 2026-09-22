<template>
  <div class="policy-description" v-html="html" />
</template>

<script setup lang="ts">
import { renderPolicyDescription } from '~/utils/markdown'

const props = withDefaults(defineProps<{ description?: string }>(), { description: '' })
const html = computed(() => renderPolicyDescription(props.description))
</script>

<style scoped>
.policy-description {
  min-width: 0;
  max-width: 100%;
  white-space: normal;
  overflow-wrap: anywhere;
}

.policy-description :deep(p),
.policy-description :deep(ul),
.policy-description :deep(ol),
.policy-description :deep(blockquote),
.policy-description :deep(pre),
.policy-description :deep(table),
.policy-description :deep(hr) {
  margin-block: 0.75em;
}

.policy-description :deep(h1),
.policy-description :deep(h2),
.policy-description :deep(h3),
.policy-description :deep(h4),
.policy-description :deep(h5),
.policy-description :deep(h6) {
  margin-block: 1em 0.5em;
  line-height: 1.3;
}

.policy-description :deep(ul),
.policy-description :deep(ol) {
  padding-inline-start: 1.5em;
}

.policy-description :deep(ul) {
  list-style-type: disc;
}

.policy-description :deep(ol) {
  list-style-type: decimal;
}

.policy-description :deep(a) {
  color: rgb(var(--v-theme-primary));
  text-decoration: underline;
}

.policy-description :deep(blockquote) {
  border-inline-start: 3px solid rgba(var(--v-theme-on-surface), 0.3);
  padding-inline-start: 1em;
}

.policy-description :deep(code) {
  background: rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 3px;
  padding: 0.15em 0.3em;
}

.policy-description :deep(pre) {
  background: rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 3px;
  padding: 0.75em;
  overflow-x: auto;
  white-space: pre;
  overflow-wrap: normal;
}

.policy-description :deep(pre code) {
  background: none;
  padding: 0;
}

.policy-description :deep(img) {
  max-width: 100%;
  height: auto;
}

.policy-description :deep(table) {
  display: block;
  max-width: 100%;
  overflow-x: auto;
  border-collapse: collapse;
}

.policy-description :deep(th),
.policy-description :deep(td) {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.3);
  padding: 0.4em 0.6em;
}

.policy-description > :deep(:first-child) {
  margin-top: 0;
}

.policy-description > :deep(:last-child) {
  margin-bottom: 0;
}
</style>
