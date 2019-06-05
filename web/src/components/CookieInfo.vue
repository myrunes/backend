<template>
  <div class="cooke-info" v-if="!hidden">
    <div class="mx-auto d-flex">
      <p class="my-auto mr-4">
        By using our site, you acknowledge and accept the use of cookies on this page.<br/>
        More information about how we are using cookies you can read about <a :href="cookielink">here</a>.
      </p>
      <button class="btn-bubble my-auto" @click="hide">OK</button>
    </div>
  </div>
</template>

<script>
import Utils from '../js/utils';

export default {
  name: 'CookieInfo',

  props: {
    cookielink: String,
  },

  data: function() {
    return {
      hidden: true,
    }
  },

  methods: {
    hide() {
      this.hidden = true;
      document.cookie = '__cookies_accepted=true; Path=/';
    }
  },

  created: function() {
    if (!Utils.getCookieValue('__cookies_accepted')) {
      this.hidden = false;
    }
  }
}
</script>

<style scoped>
  p {
    margin: 0px;
  }

  a {
    text-decoration: underline !important;
  }

  button {
    background-color: rgb(13,71,161);
    height: 40px;
    padding: 0px 20px;
  }

  .cooke-info {
    display: flex;
    position: fixed;
    bottom: 0px;
    left: 0px;
    right: 0px;
    height: fit-content;
    background-color: rgb(33,33,33);
    padding: 15px 0px;
  }

</style>
