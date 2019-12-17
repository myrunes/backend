<!-- @format -->

<template>
  <div class="outer-container">
    <div class="container my-auto">
      <div class="logo mx-auto"></div>
      <Banner
        ref="banner"
        class="mx-auto mb-5"
        width="300px"
        @closing="bannerClosing"
      ></Banner>
      <div class="d-flex position-relative">
        <b-tooltip
          ref="unameToolTipShow"
          target="tbUsername"
          triggers
          boundary="tbUsername"
          :show.sync="unameToolTipShow"
        >
          Username can only contain lower case letters, numbers, minuses (
          <code>-</code>) and underscores (
          <code>_</code>). You can later
          change your username or set a specific display name.
        </b-tooltip>
        <input
          id="tbUsername"
          v-model="username"
          type="text"
          class="tb mx-auto"
          placeholder="USERNAME"
          autocomplete="off"
          @change="checkUsername"
          @input="usernameInput"
        />
        <span class="tb mx-auto"></span>
      </div>
      <div class="d-flex mt-5 position-relative">
        <input
          id="tbPassword"
          v-model="password"
          type="password"
          class="tb mx-auto"
          placeholder="PASSWORD"
          autocomplete="off"
          @keypress="
            (e) => {
              if (e.keyCode == 13) login(e);
            }
          "
        />
        <span class="tb mx-auto"></span>
      </div>
      <div class="d-flex mt-5">
        <Slider
          v-model="remember"
          class="mx-auto"
        >
          Stay logged in (30 days)
        </Slider>
      </div>
      <div class="d-flex mt-5">
        <button
          class="btn-bubble mx-auto"
          @click="login"
        >
          {{ register ? 'REGISTER' : 'LOGIN' }}
        </button>
      </div>
      <div class="d-flex mt-5">
        <router-link
          class="text-center mx-auto forgot-password"
          to="/passwordReset"
        >
          Forgot password?
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../js/rest';
import EventBus from '../js/eventbus';
import Banner from '../components/Banner';
import Slider from '../components/Slider';

export default {
  name: 'Login',

  components: {
    Banner,
    Slider,
  },

  props: {},

  data: function() {
    return {
      register: false,

      username: '',
      password: '',
      remember: false,
      unameToolTipShow: false,
    };
  },

  mounted: function() {
    if (window.localStorage.getItem('reginfo-dismissed') !== '1') {
      this.$refs.banner.show(
        'info',
        'Not having an account? Simply register by typing in an unused username and the password you want to use.',
        null,
        true
      );
    }
  },

  methods: {
    checkUsername(e) {
      this.unameToolTipShow = false;

      let val = this.username;
      if (val.length === 0) {
        this.$refs.banner.hide();
        return;
      }

      Rest.checkUsername(val)
        .then(() => {
          this.$refs.banner.hide();
          this.register = false;
        })
        .catch((err) => {
          if (err && err.code === 404) {
            this.$refs.banner.show(
              'warning',
              'This user name is not existent. If you continue, a new account with the entered credentials will be created.',
              null,
              true
            );
            this.register = true;
          } else {
            this.$refs.banner.show(
              'error',
              `An error occured while fetching user name availability: ${
                err.message ? err.message : err
              }`,
              null,
              false
            );
          }
        });
    },

    usernameInput() {
      if (this.username.match(/([^\w_\-])|([A-Z])/g)) {
        this.unameToolTipShow = true;
      }

      this.username = this.username.toLowerCase().replace(/[^\w_\-]/g, '');
    },

    login() {
      if (!this.username || !this.password) return;

      if (this.register && this.username.length < 3) {
        this.$refs.banner.show(
          'error',
          'Username must have at least 3 characters.',
          null,
          false
        );
        return;
      }

      if (this.register && this.password.length < 8) {
        this.$refs.banner.show(
          'error',
          'Password must have at least 8 characters',
          null,
          false
        );
        return;
      }

      if (this.register) {
        Rest.register(this.username, this.password, this.remember)
          .then(() => {
            this.loginRedirect();
            window.localStorage.setItem('reginfo-dismissed', '1');
          })
          .catch((err) => {
            if (err && err.code === 409) {
              this.$refs.banner.show(
                'error',
                'This username is already in use.',
                null,
                false
              );
            } else {
              this.$refs.banner.show(
                'error',
                `An error occured during registration: ${
                  err.message ? err.message : err
                }`,
                null,
                false
              );
            }
          });
      } else {
        Rest.login(this.username, this.password, this.remember)
          .then(() => {
            this.loginRedirect();
            window.localStorage.setItem('reginfo-dismissed', '1');
          })
          .catch((err) => {
            console.log(err);
            switch (err.code) {
              case 401:
                this.$refs.banner.show(
                  'error',
                  'Invalid username-password combination.',
                  null,
                  false
                );
                break;

              case 429:
                this.$refs.banner.show(
                  'error',
                  'You have exceed your allowed ammount of login attempts. Please try again later.',
                  null,
                  false
                );
                break;

              default:
                this.$refs.banner.show(
                  'error',
                  `An error occured during login: ${
                    err.message ? err.message : err
                  }`,
                  null,
                  false
                );
            }
          });
      }
    },

    loginRedirect() {
      EventBus.$emit('login');
      let pid = this.$route.query.createpage;
      if (pid) {
        Rest.getShare(pid)
          .then((res) => {
            if (res.body && res.body.page && res.body.user) {
              res.body.page.title += ` (by ${res.body.user.displayname})`;
              Rest.createPage(res.body.page)
                .then((res) => {
                  if (res.body) {
                    this.$router.push({
                      name: 'RunePage',
                      params: { uid: res.body.uid },
                    });
                  }
                })
                .catch((err) => {
                  this.$refs.banner.show(
                    'error',
                    `Failed creating clone of shared page: ${
                      err.message ? err.message : err
                    }`,
                    null,
                    false
                  );
                });
            }
          })
          .catch((err) => {
            this.$refs.banner.show(
              'error',
              `Failed getting shared page: ${err.message ? err.message : err}`,
              null,
              false
            );
          });
      } else {
        this.$router.push('/');
      }
    },

    bannerClosing(active) {
      if (active) {
        window.localStorage.setItem('reginfo-dismissed', '1');
      }
    },
  },
};
</script>

<style scoped>
/** @format */

.outer-container {
  z-index: 100;
  width: 100vw;
  height: 100vh;
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

.forgot-password {
  text-decoration: underline !important;
  font-size: 14px;
  color: white;
  opacity: 0.8;
  transition: all 0.25s ease;
}
.forgot-password:hover {
  opacity: 1;
}
</style>
