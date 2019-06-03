/** @format */

import Router from 'vue-router';
import Login from '../components/routes/Login';
import Main from '../components/routes/Main';

export default new Router({
  mode: 'history',

  routes: [
    {
      path: '/',
      name: 'Main',
      component: Main,
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
    },
  ],
});
