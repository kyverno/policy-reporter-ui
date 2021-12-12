import Vue from 'vue'
import { mapPriority } from '~/policy-reporter-plugins/core/mapper'

Vue.filter('mapPriority', mapPriority)
