<!-- @format -->

<template>
  <div
    v-if="!hidden"
    class="cooke-info"
  >
    <div class="mx-auto d-flex">
      <p class="my-auto mr-4">
        By using our site, you acknowledge and accept the use of cookies on this
        page.<br />
        More information about how we are using cookies you can read about
        <a
          href="https://github.com/myrunes/myrunes/blob/master/docs/cookie-usage.md"
        >here</a>.
      </p>
      <button
        class="btn-bubble my-auto"
        @click="hide"
      >
        OK
      </button>
    </div>
  </div>
</template>

<script>
import Utils from '../js/utils';

export default {
  name: 'CookieInfo',

  props: {},

  data: function() {
    return {
      hidden: true,
    };
  },

  created: function() {
    if (window.localStorage.getItem('cookies-accepted') !== '1') {
      this.hidden = false;
    }
  },

  methods: {
    hide() {
      this.hidden = true;
      window.localStorage.setItem('cookies-accepted', '1');
    },
  },
};
</script>

<style scoped>
p {
  margin: 0px;
}

a {
  text-decoration: underline !important;
}

button {
  background-color: rgb(13, 71, 161);
  height: 40px;
  padding: 0px 20px;
}

.cooke-info {
  z-index: 101;
  display: flex;
  position: fixed;
  bottom: 0px;
  left: 0px;
  right: 0px;
  height: fit-content;
  background-color: rgb(33, 33, 33);
  padding: 15px 0px;
}
</style>
