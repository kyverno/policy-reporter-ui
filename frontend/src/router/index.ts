import Vue from 'vue';
import VueRouter, { RouteConfig } from 'vue-router';
import Home from '../views/Dashboard.vue';

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Dashboard',
    component: Home,
  },
  {
    path: '/policy-reports',
    name: 'PolicyReport',
    component: () => import('@/views/PolicyReport.vue'),
  },
  {
    path: '/cluster-policy-reports',
    name: 'ClusterPolicyReport',
    component: () => import('@/views/ClusterPolicyReport.vue'),
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
