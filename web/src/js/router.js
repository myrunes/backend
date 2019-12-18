/** @format */

import Router from 'vue-router';
import Login from '../routes/Login';
import Main from '../routes/Main';
import Champ from '../routes/Champ';
import Edit from '../routes/Edit';
import Pages from '../routes/Pages';
import Settings from '../routes/Settings';
import Share from '../routes/Share';
import MailConfirm from '../routes/MailConfirm';
import PasswordReset from '../routes/PasswordReset';

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
      path: '/mailConfirmation',
      name: 'MailConfirm',
      component: MailConfirm,
    },
    {
      path: '/passwordReset',
      name: 'PasswordReset',
      component: PasswordReset,
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
