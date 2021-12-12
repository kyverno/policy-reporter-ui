<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <v-card>
          <v-card-text>
            <v-alert outlined type="error" class="ma-0">
              Unable to retrieve all Data from the Server
              <v-btn v-if="action" color="error" class="ml-2" @click="action">
                reload
              </v-btn>
            </v-alert>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    action: { type: Function, default: undefined }
  }
})
</script>
