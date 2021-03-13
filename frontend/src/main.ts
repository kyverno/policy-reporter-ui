import Vue from 'vue';
import { mapPriority } from '@/models';
import VueApexCharts from 'vue-apexcharts';
import App from './App.vue';
import router from './router';
import store from './store';
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false;

Vue.filter('mapPriority', mapPriority);

Vue.use(VueApexCharts);
Vue.component('apexchart', VueApexCharts);

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),
}).$mount('#app');
