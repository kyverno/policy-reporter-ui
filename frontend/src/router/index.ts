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
  {
    path: '/logs',
    name: 'Logs',
    component: () => import('@/views/Logs.vue'),
  },
  {
    path: '/kyverno',
    component: () => import('@/plugins/kyverno/views/Dashboard.vue'),
    name: 'kyverno-dashboard',
  },
  {
    path: '/kyverno/:uid',
    name: 'policy-details',
    component: () => import('@/plugins/kyverno/views/Details.vue'),
    props: true,
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
