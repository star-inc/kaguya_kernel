/*jshint esversion: 6 */

/*
Kaguya - Web Client

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

import Login from './login.js';
import Contacts from './contacts.js';
import About from './about.js';

const routes = [{
    path: '/',
    component: Login
  },
  {
    path: '/contacts',
    component: Contacts
  },
  {
    path: '/about',
    component: About
  }
];

const router = new VueRouter({
  routes
});

const app = new Vue({
  router
}).$mount('#app');
