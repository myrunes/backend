<template>
  <div>
    <Banner v-if="banner.visible" :type="banner.type" class="mb-3">{{ banner.content }}</Banner>
    <div class="bg mb-3">
      <h3>ACCOUNT DETAILS</h3>
      <table>
        <tr>
          <td class="pr-5">Created</td>
          <td>{{ formatTime(user.created) }}</td>
        </tr>
        <tr>
          <td class="pr-5">Last Login</td>
          <td>{{ formatTime(user.lastlogin) }}</td>
        </tr>
        <tr>
          <td class="pr-5">Pages</td>
          <td>{{ pages }}</td>
        </tr>
        <tr>
          <td class="pr-5">Created</td>
          <td class="hider">{{ user.uid }}</td>
        </tr>
      </table>
    </div>
    <div class="bg">
      <h3 class="mb-3">UPDATE ACCOUNT</h3>

      <div class="position-relative mb-4">
        <h5>Username</h5>
        <p class="explainer">
          The unique identifyer you need to use to log in.<br/>
          The username must be lowercase, longer than 3 characters and must only contain letters, numbers, scores and underscores.
        </p>
        <input type="text" class="tb text-left" v-model="user.username" @input="unameInput"/>
        <span class="tb"/>
      </div>

      <div class="position-relative mb-4">
        <h5>Display Name</h5>
        <p class="explainer">
          The name which may be displayed to other users.
        </p>
        <input type="text" class="tb text-left" v-model="user.displayname"/>
        <span class="tb"/>
      </div>

      <div class="position-relative">
        <h5>New Password</h5>
        <p class="explainer">
          Enter a new password, if you want to change it.
        </p>
        <input type="password" ref="tbNewpw" class="tb text-left" v-model="newpassword"/>
        <span class="tb"/>
        <a class="ml-2" @mousedown="$refs.tbNewpw.type='text'" @mouseup="$refs.tbNewpw.type='password'">
          <img src="/assets/eye.svg" width="20" height="20"/>
        </a>
      </div>

      <div class="mt-5">
        <hr>
        <p>You need to enter your current password again to apply these changes:</p>
        <div class="position-relative mb-4">
          <input type="password" class="tb text-left" v-model="currpassword"/>
          <span class="tb"/>
        </div>
        <div class="bg danger-zone mb-3">
          <h5 class="mb-3">DANGER ZONE</h5>
          <button class="btn-slide btn-delete" @click="deleteAcc">DELETE ACCOUNT PERMANENTLY AND FOREVER</button>
        </div>
        <div class="text-right">
          <button class="btn-slide btn-save mr-3" @click="save">SAVE</button>
          <button class="btn-slide btn-cancel" @click="$router.back()">CANCEL</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Rest from '../../js/rest';
import Banner from '../Banner';
import Utils from '../../js/utils';
import EventBus from '../../js/eventbus';

export default {
  name: 'Settings',

  components: {
    Banner,
  },

  data: function() {
    return {
      user: {},
      pages: 0,
      newpassword: '',
      currpassword: '',
      
      banner: {
        visible: false,
        type: 'error',
        content: '',
      },
    }
  },

  methods: {
    formatTime(timestamp) {
      let t = Utils.parseTime(new Date(timestamp));
      return `${t.y}/${t.m}/${t.d} ${t.h}:${t.m}:${t.s}`;
    },

    unameInput(e) {
      this.user.username = this.user.username
        .toLowerCase()
        .replace(/[^\w_\-]/g, '');
    },

    save() {
      let currpw = this.currpassword;
      this.currpassword = '';

      if (this.newpassword && this.newpassword.length < 8) {
        this.banner = {
          visible: true,
          type: 'error',
          content: `Password must have at least 8 characters.`,
        }
        setTimeout(() => this.banner.visible = false, 10000);
        window.scrollTo(0, 0);
        return;
      } 

      let update = {
        username: this.user.username,
        displayname: this.user.displayname,
        currpassword: currpw,
        newpassword: this.newpassword,
      }
      Rest.updateUser(update).then(() => {
        this.banner = {
          visible: true,
          type: 'success',
          content: `Account changes saved.`,
        }
        setTimeout(() => this.banner.visible = false, 10000);
        window.scrollTo(0, 0);
      }).catch((err) => {
        this.banner = {
          visible: true,
          type: 'error',
          content: err.message === 'unauthorized' ? 'Current password is wrong.' : `Saving failed: ${err.message ? err.message : err}`,
        }
        setTimeout(() => this.banner.visible = false, 10000);
        window.scrollTo(0, 0);
        console.error(err);
      });
    },

    deleteAcc() {
      let currpw = this.currpassword;
      this.currpassword = '';

      Rest.deleteUser(currpw).then(() => {
        EventBus.$emit('logout');
        this.$router.push('/login');
      }).catch((err) => {
        this.banner = {
          visible: true,
          type: 'error',
          content: err.message === 'unauthorized' ? 'Current password is wrong.' : `Saving failed: ${err.message ? err.message : err}`,
        }
        setTimeout(() => this.banner.visible = false, 10000);
        window.scrollTo(0, 0);
        console.error(err);
      });
    }
  },

  created: function() {
    Rest.getMe().then((res) => {
      if (!res.body) return;
      this.user = res.body;
      console.log(this.user);
    }).catch(console.error);

    Rest.getPages().then((res) => {
      if (!res.body) return;
      this.pages = res.body.n;
    }).catch(console.error)
  },
}

</script>


<style scoped>

.hider {
  color: rgb(33,33,33);
  background-color: rgb(33,33,33);
}

.hider:hover {
  color: white;
  background-color: transparent;
}

.hider:hover::before {
  display: none;
}

.hider::before {
  content: 'HOVER TO DISPLAY';
  color: rgb(133, 133, 133);
  font-size: 14px;
  position: absolute;
  transform: translate(4px, 1px);
}

h5, p {
  margin: 0px;
}

h5 {
  font-family: 'Montserrat', sans-serif;
}

.explainer {
  font-style: italic;
  font-size: 14px;
}

.danger-zone {
  background-color: rgba(198, 40, 40, 0.3);
  width: fit-content;
  height: fit-content;
}

.btn-delete {
  background-color: rgb(244,67,54);
}

.btn-delete::before {
  background-color: rgb(213,0,0);
}

</style>
