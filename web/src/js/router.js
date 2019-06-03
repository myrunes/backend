/** @format */

import Router from 'vue-router';
import Login from '../components/routes/Login';

export default new Router({
  mode: 'history',

  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login,
    },
  ],
});
