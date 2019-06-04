<template>
  <div id="app">
    <Header v-if="loggedIn" />
    <router-view :class="{ m : loggedIn }"></router-view>
  </div>
</template>

<script>
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';
import './css/global.css';

import Rest from './js/rest';
import Router from './js/router';
import EventBus from './js/eventbus';

import Header from './components/Header';

export default {
  name: 'app',

  components: {
    Header,
  },

  router: Router,

  data: function() {
    return {
      loggedIn: false,
    }
  },

  created: function() {
    this.checkLogin();

    EventBus.$on('login', () => {
      this.checkLogin();
    });

    EventBus.$on('logout', () => {
      this.loggedIn = false;
    });
  },

  methods: {
    checkLogin() {
      Rest.getMe().then((res) => {
        this.loggedIn = true;
      }).catch((err) => {
        this.$router.push('/login');
      })
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.m {
  margin: 70px 20px;
}
</style>
