<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <v-card>
          <v-card-text>
            <h1 v-if="error.statusCode === 404">
              {{ pageNotFound }}
            </h1>
            <h1 v-else>
              {{ otherError }}
            </h1>
          </v-card-text>
          <v-divider />
          <v-card-actions>
            <v-btn outlined to="/" color="error">
              Back to Home
            </v-btn>
            <v-btn outlined color="error" @click="refresh">
              Reload the Page
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  layout: 'empty',
  props: {
    error: {
      type: Object,
      default: null
    }
  },
  data () {
    return {
      pageNotFound: '404 Not Found',
      otherError: 'Configured backend not accessible'
    }
  },
  head () {
    const title =
      this.error.statusCode === 404 ? this.pageNotFound : this.otherError
    return {
      title
    }
  },
  methods: {
    refresh () {
      location.reload()
    }
  }
}
</script>

<style scoped>
h1 {
  font-size: 20px;
}
</style>
