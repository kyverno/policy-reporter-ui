<template>
    <v-autocomplete dense
                    multiple
                    :items="policies"
                    outlined
                    hide-details
                    label="Policies"
                    :value="value"
                    v-bind="$attrs"
                    v-on="$listeners"
    >
    <template v-slot:selection="{ item, index }">
        <v-chip small v-if="index <= 1" label outlined>
          <span>{{ item }}</span>
        </v-chip>
        <span
          v-if="index === 2"
          class="grey--text caption"
        >
          (+{{ value.length - 2 }} others)
        </span>
      </template>
    </v-autocomplete>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend<{}, {}, { policies: string[] }, { policies: string[]; value: string[] }>({
  props: {
    policies: { type: Array, default: () => [] },
    value: { type: Array, default: () => [] },
  },
  created() {
    if (!this.policies.length || this.value.length > 0) return;

    this.$emit('input', [this.policies[0]]);
  },
});
</script>
