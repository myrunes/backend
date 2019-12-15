<template>
  <div id="app">
    <CookieInfo />
    <Header v-if="loggedIn" />
    <router-view :class="{ m : loggedIn }"></router-view>
    <InfoBubble ref="betawarn" color="red" @hides="onBetaWarnHides" style="z-index: 110;">
      <p class="text-center m-2">
        <b>ATTENTION:</b>&nbsp;This instance is running on a non-release canary build which is not fully tested yet!
        <br />Actions taken here may lead to loss or corruption of data!
      </p>
    </InfoBubble>
    <Footer />
  </div>
</template>

<script>
/** @format */

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';
import './css/global.css';

import Rest from './js/rest';
import Router from './js/router';
import EventBus from './js/eventbus';

import Header from './components/Header';
import Footer from './components/Footer';
import CookieInfo from './components/CookieInfo';
import InfoBubble from './components/InfoBubble';

const NO_LOGIN_ROUTES = ['Share', 'MailConfirm', 'PasswordReset'];

export default {
  name: 'app',

  components: {
    Header,
    CookieInfo,
    Footer,
    InfoBubble,
  },

  router: Router,

  data: function() {
    return {
      loggedIn: false,
    };
  },

  created: function() {
    this.checkLogin();

    Rest.getVersion().then((res) => {
      if (
        res.body.release !== 'TRUE' &&
        window.localStorage.getItem('beta-info-accepted') !== '1'
      ) {
        setTimeout(() => this.$refs.betawarn.show(), 1000);
      }
    });

    EventBus.$on('login', () => {
      this.checkLogin();
    });

    EventBus.$on('logout', () => {
      this.loggedIn = false;
    });
  },

  methods: {
    checkLogin() {
      Rest.getMe()
        .then((res) => {
          this.loggedIn = true;
        })
        .catch((err) => {
          if (!NO_LOGIN_ROUTES.includes(this.$route.name)) {
            this.$router.replace('/login');
          }
        });
    },

    onBetaWarnHides() {
      window.localStorage.setItem('beta-info-accepted', '1');
    },
  },
};
</script>

<style>
/** @format */

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.m {
  margin: 70px 20px;
}
</style>
