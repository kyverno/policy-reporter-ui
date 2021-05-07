import Vue from 'vue';
import { mapPriority } from '@/mapper';
import VueApexCharts from 'vue-apexcharts';
import hljs from 'highlight.js';
import VueClipboard from 'vue-clipboard2';
import App from './App.vue';
import router from './router';
import store from './store';
import vuetify from './vuetify';
import 'highlight.js/styles/github.css';

Vue.config.productionTip = false;

Vue.filter('mapPriority', mapPriority);

Vue.use(hljs.vuePlugin);
Vue.use(VueApexCharts);
Vue.component('apexchart', VueApexCharts);

VueClipboard.config.autoSetContainer = true;
Vue.use(VueClipboard);

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),
}).$mount('#app');
