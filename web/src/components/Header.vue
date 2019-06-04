<template>
  <div class="header d-flex">
    <div class="d-flex">
      <router-link to="/">
        <button class="btn-slide btn-logo">
          <h1 class="mr-4">RUNES.WTF</h1>
        </button>
      </router-link>
      <button class="btn-slide">MY PAGES</button>
      <button class="btn-slide">SETTINGS</button>
    </div>
    <button class="btn-slide btn-logout" @click="logout">LOGOUT</button>
  </div>
</template>

<script>
import Rest from '../js/rest';
import EventBus from '../js/eventbus';

export default {
  name: 'Header',

  props: {
  },

  methods: {
    logout() {
      Rest.logout().then(() => {
        EventBus.$emit('logout');
        this.$router.push('/login');
      }).catch(console.error);
    }
  }  
}
</script>

<style>

.header > button {
  position: relative;
}

.header {
  position: fixed;
  top: 0px;
  left: 0px;
  right: 0px;
  height: 50px;
  background-color: #37474F;
  justify-content: space-between;
}

.btn-logo {
  background-color: transparent;
}

.btn-logo::before {
  background-color: #03A9F4;
}

.btn-logout::before {
  background-color: #FF5722;
}

h1 {
  font-size: 26px;
  margin-bottom: auto;
  margin-top: auto;
  margin-left: 20px;
}

</style>
