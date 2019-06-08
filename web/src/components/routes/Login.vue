<template>
  <div class="outer-container">
    <div class="container my-auto">
      <div class="logo mx-auto"></div>
      <Banner 
        class="mx-auto mb-5"
        width="300px"
        v-if="banner.visible" 
        :type="banner.type"
        :closable="banner.closable"
        @closing="bannerClosing('reginfo')"
      >
        {{ banner.content }}
      </Banner>
      <div class="d-flex position-relative">
        <input 
          id="tbUsername" 
          type="text" 
          class="tb mx-auto"
          @change="checkUsername"
          v-model="username"
          placeholder="USERNAME"
          autocomplete="off"
        >
        <span class="tb mx-auto"></span>
      </div>
      <div class="d-flex mt-5 position-relative">
        <input 
          v-model="password" 
          id="tbPassword" 
          type="password" 
          class="tb mx-auto"
          placeholder="PASSWORD"
          autocomplete="off"
          @keypress="(e) => { if (e.keyCode == 13) login(e) }"
        >
        <span class="tb mx-auto"></span>
      </div>
      <div class="d-flex mt-5">
        <Slider 
          class="mx-auto"
          v-model="remember"
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
    </div>
  </div>
</template>

<script>
import Rest from '../../js/rest';
import Banner from '../Banner';
import Slider from '../Slider';
import EventBus from '../../js/eventbus';

export default {
  name: 'Login',
  
  props: {
  },

  components: {
    Banner,
    Slider,
  },

  data: function() {
    return {
      banner: {
        visible: false,
        type: 'warning',
        content: ''
      },

      register: false,

      username: '',
      password: '',
      remember: false,
    }
  },

  methods: {
    checkUsername(e) {
      let val = e.target.value;
      if (val.length === 0) {
        this.banner.visible = false;
        return;
      }

      Rest.checkUsername(val).then(() => {
        this.banner.visible = false;
        this.register = false;
      }).catch((err) => {
        if (err && err.code === 404) {
          this.banner = {
            visible: true,
            type: 'warning',
            content: `This user name is not existent. If you continue, a new account with the entered credentials will be created.`
          }
          this.register = true;
        } else {
          this.banner = {
            visible: true,
            type: 'error',
            cintent: `An error occured while fetching user name availability: ${err.message ? err.message : err}`
          }
        }
      });
    },

    login() {
      if (!this.username || !this.password) return;

      if (this.register && this.username.length < 3) {
        this.banner = {
          visible: true,
          type: 'error',
          content: 'Username must have at least 3 characters.'
        }
        return;
      }

      if (this.register && this.password.length < 8) {
        this.banner = {
          visible: true,
          type: 'error',
          content: 'Password must have at least 8 characters.'
        }
        return;
      }

      if (this.register) {
        Rest.register(this.username, this.password, this.remember).then(() => {
          EventBus.$emit('login');
          this.$router.push('/');
          window.localStorage.setItem('reginfo-dismissed', '1');
        }).catch((err) => {
          if (err && err.code === 409) {
            this.banner = {
              visible: true,
              type: 'error',
              content: 'The passes username is already in use.'
            }
          } else {
            this.banner = {
              visible: true,
              type: 'error',
              content: `An error occured during registration: ${err.message ? err.message : err}`
            }
          }
        });
      } else {
        Rest.login(this.username, this.password, this.remember).then(() => {
          EventBus.$emit('login');
          this.$router.push('/');
        }).catch((err) => {
          if (err && err.code === 401) {
            this.banner = {
              visible: true,
              type: 'error',
              content: 'Invalid username-password combination.'
            }
          } else {
            this.banner = {
              visible: true,
              type: 'error',
              content: `An error occured during registration: ${err.message ? err.message : err}`
            }
          }
        });
      }
      
    },

    bannerClosing(reason) {
      this.banner.visible = false;
      if (reason === 'reginfo') {
        window.localStorage.setItem('reginfo-dismissed', '1');
      }
    }
  },

  created: function() {
    if (window.localStorage.getItem('reginfo-dismissed') !== '1') {
      this.banner = {
        closable: true,
        visible: true,
        type: 'info',
        content: 'Not having an account? Simply register by typing in an unused username and the password you want to use.'
      }
    }
  }

}
</script>

<style scoped>

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

</style>
