<template>
  <app-row>
    <v-card>
      <v-toolbar color="category">
        <v-toolbar-title>Policy Details</v-toolbar-title>
        <template #append>
          <CollapseBtn v-model="show"/>
        </template>
      </v-toolbar>
    </v-card>
  </app-row>

  <v-expand-transition>
    <v-row v-if="show">
      <v-col v-if="policy.engine" cols="12" lg="5">
        <v-card title="Engine Information">
          <v-container>
            <v-row class="top-border">
              <v-col class="font-weight-bold">Name</v-col>
              <v-col>{{ policy.engine.name }}</v-col>
            </v-row>
            <v-row class="top-border" v-if="policy.engine.version">
              <v-col class="font-weight-bold">Min. Version</v-col>
              <v-col>{{ policy.engine.version }}</v-col>
            </v-row>
            <v-row class="top-border" v-if="(policy.engine.subjects || []).length">
              <v-col class="font-weight-bold">Subjects</v-col>
              <v-col>{{ policy.engine.subjects.join(', ') }}</v-col>
            </v-row>
          </v-container>
        </v-card>
        <v-card v-if="policy.details && policy.details.length" class="mt-4" title="Details">
          <v-container>
            <v-row class="top-border" v-for="item in policy.details" :key="item.value">
              <v-col class="font-weight-bold" v-if="item.title">{{ item.title }}</v-col>
              <v-col>{{ item.value }}</v-col>
            </v-row>
          </v-container>
        </v-card>
      </v-col>

      <v-col cols="12" lg="7">
        <v-card style="height: 100%" title="Description">
          <v-divider/>
          <v-card-text class="text-pre-line overflow-y-hidden" style="height: calc(100% - 104px)">
            <v-expand-transition>
            <div :style="{ maxHeight: descriptionHeight }">
              {{ policy.description }}
            </div>
            </v-expand-transition>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn @click="expand = !expand">expand</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>

      <template v-if="policy.additional && policy.additional.length">
        <v-col v-for="details in policy.additional" :cols="detailsCols">
          <v-card :title="details.title">
            <v-container fluid>
              <template v-for="item in details.items" :key="item.value">
                <v-row class="top-border" v-if="item.value">
                  <v-col :cols="detailsLabelCols.label" class="font-weight-bold" v-if="item.title">{{ item.title }}</v-col>
                  <v-col :cols="detailsLabelCols.value">{{ item.value }}</v-col>
                </v-row>
              </template>
            </v-container>
          </v-card>
        </v-col>
      </template>
    </v-row>
  </v-expand-transition>

  <app-row v-if="policy.sourceCode && policy.sourceCode.content">
    <v-card>
      <v-toolbar color="category">
        <v-toolbar-title>Source Code</v-toolbar-title>
        <template #append>
          <CollapseBtn v-model="open"/>
        </template>
      </v-toolbar>
      <v-expand-transition>
        <div v-if="open">
          <v-divider/>
          <v-card-text>
            <highlightjs :code="policy.sourceCode.content" :lang="policy.sourceCode.contentType"/>
          </v-card-text>
          <v-divider/>
          <v-card-actions>
            <v-spacer/>
            <v-btn @click="open = false">Close</v-btn>
          </v-card-actions>
        </div>
      </v-expand-transition>
    </v-card>
  </app-row>

  <app-row v-if="policy.references && policy.references.length">
    <v-card>
      <v-toolbar color="category">
        <v-toolbar-title>References</v-toolbar-title>
        <template #append>
          <CollapseBtn v-model="references"/>
        </template>
      </v-toolbar>
      <v-expand-transition>
        <v-list class="pt-0 mt0" lines="one" v-if="references">
          <template v-for="link in policy.references" :key="link">
            <v-divider />
            <v-list-item>
              <v-list-item-title>
                <a :href="link" class="text-primary text-decoration-none" target="_blank">{{ link }}</a>
              </v-list-item-title>
            </v-list-item>
          </template>
        </v-list>
      </v-expand-transition>
    </v-card>
  </app-row>
</template>

<script lang="ts" setup>
import { type PolicyDetails } from "../../types";

const props = defineProps<{ policy: PolicyDetails; defaultOpen?: boolean; show?: boolean }>()

const open = ref(props.defaultOpen ?? false)
const show = ref(true)
const references = ref(false)
const expand = ref(false)

const descriptionHeight = computed(() => expand.value ? undefined : '200px')

const detailsCols = computed(() => {
  return Math.ceil(12 / props.policy.additional.length)
})

const detailsLabelCols = computed(() => {
  if (detailsCols.value === 12) return ({ label: 2, value: 10 })

  return { label: 6, value: 6 }
})
</script>

<style scoped>

</style>
