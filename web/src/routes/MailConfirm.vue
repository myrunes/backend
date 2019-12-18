<!-- @format -->

<template>
  <div class="outer-container">
    <div class="container my-auto text-center">
      <div class="logo mx-auto"></div>
      <Banner
        ref="banner"
        class="mx-auto mb-5"
      ></Banner>
      <h2>Mail Address Confirmation</h2>
      <p class="mt-4">
        Click the following button to confirm your Mail Address.
      </p>
      <button
        class="btn-bubble mx-auto mt-4"
        @click="confirmClick"
      >
        Confirm Mail Address
      </button>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../js/rest';
import Banner from '../components/Banner';

export default {
  name: 'Login',

  components: {
    Banner,
  },

  props: {},

  data: function() {
    return {
      token: '',
    };
  },

  mounted: function() {
    this.token = this.$route.query.token;
  },

  methods: {
    confirmClick() {
      Rest.confirmMail(this.token)
        .then(() => {
          this.$router.replace('/settings');
        })
        .catch((err) => {
          this.$refs.banner.show(
            'error',
            'The passed confirmation token is invalid!',
            null,
            true
          );
          console.error(err);
        });
    },
  },
};
</script>

<style scoped>
/** @format */

.outer-container {
  margin-top: 20vh;
  z-index: 100;
  display: flex;
}

button {
  width: 300px;
}

.logo {
  width: 256px;
  height: 61px;
  background-image: url('/assets/logo-256-61.png');
  background-repeat: no-repeat;
  background-size: 100%;
  background-position: center;
  margin-bottom: 50px;
}
</style>
