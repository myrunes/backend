/** @format */

import Router from 'vue-router';
import Login from '../components/routes/Login';
import Main from '../components/routes/Main';
import Champ from '../components/routes/Champ';
import Edit from '../components/routes/Edit';
import Pages from '../components/routes/Pages';
import Settings from '../components/routes/Settings';
import Share from '../components/routes/Share';

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
    {
      path: '/champ/:champ',
      name: 'Champ',
      component: Champ,
    },
    {
      path: '/page/:uid',
      name: 'RunePage',
      component: Edit,
    },
    {
      path: '/pages',
      name: 'RunePages',
      component: Pages,
    },
    {
      path: '/settings',
      name: 'Settings',
      component: Settings,
    },
    {
      path: '/p/:ident',
      name: 'Share',
      component: Share,
    },
  ],
});
