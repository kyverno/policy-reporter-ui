<template>
  <v-select
    :value="currentCluster"
    :items="clusters"
    outlined
    hide-details
    dense
    item-text="name"
    item-value="id"
    style="max-width: 200px;"
    @input="select"
  />
</template>

<script lang="ts">
import Vue from 'vue'
import { mapGetters, mapMutations } from 'vuex'
import { Cluster } from '~/policy-reporter-plugins/core/types'

type Data = {}

type Methdos = {
  setCluster(cluster?: Cluster): void;
  select(id: string): void
}

type Computed = {
  currentCluster?: Cluster;
  clusters: Cluster[]
}

export default Vue.extend<Data, Methdos, Computed, {}>({
  computed: mapGetters(['clusters', 'currentCluster']),
  methods: {
    ...mapMutations(['setCluster']),
    select (id: string): void {
      const cluster = this.clusters.find(c => c.id === id)

      if (!cluster || !cluster.id) {
        this.$coreAPI.setPrefix('')
        this.$kyvernoAPI.setPrefix('')
        this.setCluster(cluster)
        return
      }

      this.$coreAPI.setPrefix('/' + cluster.id)

      if (cluster.kyverno) {
        this.$kyvernoAPI.setPrefix('/' + cluster.id)
      }

      this.setCluster(cluster)
    }
  }
})
</script>
