<!-- @format -->

<template>
  <div class="header d-flex">
    <div class="d-flex">
      <router-link to="/">
        <button class="btn-slide btn-logo">
          <img src="/assets/logo-256-61.png" />
        </button>
      </router-link>
      <router-link to="/pages">
        <button class="btn-slide h-100">
          MY PAGES
        </button>
      </router-link>
      <router-link to="/settings">
        <button class="btn-slide h-100">
          MY SETTINGS
        </button>
      </router-link>
    </div>
    <button
      class="btn-slide btn-logout"
      @click="logout"
    >
      LOGOUT
    </button>
  </div>
</template>

<script>
import Rest from '../js/rest';
import EventBus from '../js/eventbus';

export default {
  name: 'Header',

  props: {},

  methods: {
    logout() {
      Rest.logout()
        .then(() => {
          EventBus.$emit('logout');
          this.$router.push('/login');
        })
        .catch(console.error);
    },
  },
};
</script>

<style scoped>
.header > button {
  position: relative;
}

.header {
  position: fixed;
  z-index: 99;
  top: 0px;
  left: 0px;
  right: 0px;
  height: 50px;
  background-color: #37474f;
  justify-content: space-between;
}

.btn-logo {
  /* background-color: transparent; */
  padding: 10px;
  height: 100%;
  width: 180px;
}

.btn-logo > img {
  height: 100%;
}

.btn-logo::before {
  background-color: #03a9f4;
}

.btn-logout::before {
  background-color: #ff5722;
}

h1 {
  font-size: 26px;
  margin-bottom: auto;
  margin-top: auto;
  margin-left: 20px;
}

b {
  font-weight: 400;
}
</style>
