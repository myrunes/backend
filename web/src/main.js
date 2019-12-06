/** @format */

import Vue from 'vue';
import App from './App.vue';
import BootstrapVue from 'bootstrap-vue';
import Router from 'vue-router';
import AsyncComputed from 'vue-async-computed';

import { Sortable, AutoScroll } from 'vuedraggable';

Vue.config.productionTip = false;

Vue.use(BootstrapVue);
Vue.use(Router);
Vue.use(AsyncComputed);

new Vue({
  render: (h) => h(App),
}).$mount('#app');
