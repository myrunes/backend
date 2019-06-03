<template>
  <div class="header d-flex">
    <div class="d-flex">
      <h1 class="mr-4">LOL RUNES</h1>
      <button class="btn-slide">MY PAGES</button>
      <button class="btn-slide">SETTINGS</button>
    </div>
    <button class="btn-slide" @click="logout">LOGOUT</button>
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

h1 {
  font-size: 26px;
  margin-bottom: auto;
  margin-top: auto;
  margin-left: 20px;
}

</style>
