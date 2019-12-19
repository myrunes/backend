<!-- @format -->

<template>
  <div class="outer-container">
    <div class="mx-auto text-center">
      <div class="logo mx-auto"></div>
      <Banner ref="banner" class="mx-auto mb-5"></Banner>

      <!-- PASSWORD RESET FORM -->
      <div v-if="token != null && token.length === 0" class="mx-auto max-w-container">
        <h2>Password Reset</h2>
        <p class="mt-5">Please enter your mail address to reset your password.</p>
        <p>Then, we will send you an E-Mail with the reset confirmation.</p>
        <input
          v-model="mailaddress"
          type="text"
          class="tb text-center mt-4"
          placeholder="Mail Address"
        />
        <br />
        <button
          class="btn-bubble mx-auto mt-5"
          :disabled="mailDisabled()"
          @click="resetClick"
        >RESET PASSWORD</button>
      </div>

      <!-- SSECURITY CHECK -->
      <div v-if="token != null && token.length > 0" class="mx-auto mb-5 max-w-container">
        <h2>Password Reset Security Check</h2>
        <p class="mt-4">Please enter your new password</p>
        <b-tooltip
          target="passwordInput"
          triggers
          boundary="passwordInput"
          :show="password.length > 0 && password.length < 8"
        >The password must have at least a length of 8 characters.</b-tooltip>
        <input
          id="passwordInput"
          v-model="password"
          type="password"
          class="tb text-center"
          placeholder="New Password"
        />
        <br />
        <input
          v-model="passwordRepeated"
          type="password"
          class="tb text-center mt-4"
          placeholder="Repeat Password"
        />
        <p
          class="mt-5"
        >To ensure that you are really the owner of this account, please enter 3 names of rune pages you have created.</p>
        <p class="smal-text">
          If you have less than 3 pages or if you can not remember their names and you
          really need your account back, please contact us via our contact mail address.
        </p>
        <input v-model="page_names[0]" type="text" class="tb text-center" placeholder="Page Name 1" />
        <br />
        <input
          v-model="page_names[1]"
          type="text"
          class="tb text-center mt-4"
          placeholder="Page Name 2"
        />
        <br />
        <input
          v-model="page_names[2]"
          type="text"
          class="tb text-center mt-4"
          placeholder="Page Name 3"
        />
        <br />
        <button
          class="btn-bubble mx-auto mt-5"
          :disabled="confirmDisabled()"
          @click="resetConfirmClick"
        >SET PASSWORD</button>
      </div>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../js/rest';
import EventBus from '../js/eventbus';
import Banner from '../components/Banner';

export default {
  name: 'Login',

  components: {
    Banner,
  },

  props: {},

  data: function() {
    return {
      passwordToolTipShow: false,

      token: null,

      mailaddress: '',
      password: '',
      passwordRepeated: '',
      page_names: [],
    };
  },

  mounted: function() {
    this.token = this.$route.query.token || '';
  },

  methods: {
    resetClick() {
      Rest.resetPassword(this.mailaddress)
        .then(() => {
          this.mailaddress = '';
          this.$refs.banner.show(
            'info',
            'We have sent an E-Mail to the address to reset your password.\nPlease also check your spam folder for this mail.',
            null,
            true
          );
        })
        .catch(console.error);
    },

    resetConfirmClick() {
      Rest.resetPasswordConfirm(this.token, this.password, this.page_names)
        .then(() => {
          this.password = '';
          this.passwordRepeated = '';
          this.page_names = [];
          this.$refs.banner.show(
            'info',
            'Your password was reset. You can now log in using your new password.',
            null,
            true
          );
        })
        .catch((err) => {
          let txt = 'Password reset failed.';
          if (err.message) {
            txt = `Password reset failed: ${err.message}`;
          }
          this.$refs.banner.show('error', txt, null, true);
        });
    },

    mailDisabled() {
      return !this.mailaddress || !this.mailaddress.includes('@');
    },

    confirmDisabled() {
      return (
        !this.password ||
        !this.passwordRepeated ||
        this.password.length < 8 ||
        this.password !== this.passwordRepeated ||
        !this.page_names[0] ||
        !this.page_names[1] ||
        !this.page_names[2]
      );
    },
  },
};
</script>

<style scoped>
/** @format */

.outer-container {
  margin-top: 15vh;
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

.max-w-container {
  max-width: 600px;
}

.smal-text {
  font-size: 14px;
  color: rgb(185, 185, 185);
}
</style>
