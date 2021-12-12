<template>
  <v-card style="height: 100%;">
    <v-toolbar flat>
      <v-toolbar-title>Details</v-toolbar-title>
    </v-toolbar>
    <v-list>
      <v-divider />
      <v-list-item two-line>
        <v-list-item-content>
          <v-list-item-title>
            {{ policy.category || 'No Category' }}
          </v-list-item-title>
          <v-list-item-subtitle>Category</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list-item two-line>
        <v-list-item-content>
          <v-list-item-title>
            {{ policy.namespace || 'Cluster scoped' }}
          </v-list-item-title>
          <v-list-item-subtitle>Namespace</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-divider v-if="policy.severity" />
      <v-list-item v-if="policy.severity" two-line>
        <v-list-item-content>
          <v-list-item-title>
            {{ policy.severity || 'Cluster scoped' }}
          </v-list-item-title>
          <v-list-item-subtitle>Severity</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-divider v-if="policy.validationFailureAction" />
      <v-list-item v-if="policy.validationFailureAction" two-line>
        <v-list-item-content>
          <v-list-item-title>
            {{ policy.validationFailureAction }}
          </v-list-item-title>
          <v-list-item-subtitle>Validation Failure Action</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { Policy } from '../types'

export default Vue.extend<{}, {}, {}, { policy: Policy }>({
  props: {
    policy: { type: Object, required: true }
  }
})
</script>
