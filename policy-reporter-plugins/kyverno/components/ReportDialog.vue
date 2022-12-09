<template>
  <v-dialog v-model="open" width="600" persistent>
    <template #activator="{ on, attrs }">
      <v-btn
        dark
        class="action-button"
        v-bind="attrs"
        v-on="on"
      >
        Generate
      </v-btn>
    </template>

    <v-card>
      <v-card-title>
        Generate Compliance Report
        <v-spacer />
        <v-btn icon @click="open = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-card-text>
        <v-container>
          <v-row>
            <v-col>
              <policy-autocomplete v-model="policies" all source="kyverno" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <namespace-autocomplete v-model="namespaces" source="kyverno" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <repor-type-select v-model="type" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-switch v-model="clusterScope" label="include cluster scoped results" class="mt-0" />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              <v-btn color="error" @click="open = false">
                Close
              </v-btn>
            </v-col>
            <v-col cols="4" offset="1">
              <report-button
                block
                :type="type"
                :policies="policies"
                :namespaces="namespaces"
                :cluster-scope="clusterScope"
                open
              />
            </v-col>
            <v-col cols="4">
              <report-button block :type="type" :policies="policies" :namespaces="namespaces" :cluster-scope="clusterScope" />
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import Vue from 'vue'
import { ReportType } from '../types'
import ReporTypeSelect from './ReporTypeSelect.vue'
import ReportButton from './ReportButton.vue'
import PolicyAutocomplete from '~/components/selects/PolicyAutocomplete.vue'
import NamespaceAutocomplete from '~/components/selects/NamespaceAutocomplete.vue'

type Data = {
  open: boolean;
  policies: string[];
  namespaces: string[];
  type: ReportType;
  clusterScope: boolean;
};
type Computed = {};
type Props = {};
type Methods = {};

export default Vue.extend<Data, Methods, Computed, Props>({
  components: { PolicyAutocomplete, NamespaceAutocomplete, ReporTypeSelect, ReportButton },
  props: {},
  data: () => ({
    open: false,
    namespaces: [],
    policies: [],
    clusterScope: true,
    type: ReportType.POLICY
  }),
  watch: {
    open (open: boolean) {
      if (!open) {
        setTimeout(() => {
          this.policies = []
          this.namespaces = []
          this.clusterScope = true
          this.type = ReportType.POLICY
        }, 200)
      }
    }
  }
})
</script>

<style scoped>
>>> code {
  padding: 16px !important;
}
</style>
