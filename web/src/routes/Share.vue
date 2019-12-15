<!-- @format -->

<template>
  <div
    class="outer-container"
    :class="{ 'm-3': !loggedin }"
  >
    <div
      v-if="notfound"
      class="container-fluid d-flex"
    >
      <h1 class="mx-auto mt-5">
        This is not a shared page.&nbsp;&nbsp;:(
      </h1>
    </div>

    <div v-else>
      <div class="d-flex">
        <h1>{{ page.title }}</h1>
        <div class="d-flex my-auto ml-4">
          <img
            v-for="c in page.champions"
            :key="`champ-${c}`"
            class="round mr-3"
            height="30"
            width="30"
            :src="`/assets/champ-avis/${c}.png`"
          />
        </div>
      </div>
      <h4>by {{ user.displayname }}</h4>
      <div
        class="d-flex mt-4 bg-runes"
        :class="`tree-${page.primary.tree}`"
      >
        <img
          :src="`/assets/rune-avis/${page.primary.tree}.png`"
          width="70"
          height="70"
          class="mr-3"
        />
        <img
          v-for="r in page.primary.rows"
          :key="`rune-${r}`"
          :src="`/assets/rune-avis/${page.primary.tree}/${r}.png`"
          width="70"
          height="70"
          class="mr-2"
        />
      </div>
      <div
        class="d-flex mt-4 bg-runes"
        :class="`tree-${page.secondary.tree}`"
      >
        <img
          :src="`/assets/rune-avis/${page.secondary.tree}.png`"
          width="70"
          height="70"
          class="mr-3"
        />
        <img
          v-for="r in page.secondary.rows"
          :key="`rune-${r}`"
          :src="`/assets/rune-avis/${page.secondary.tree}/${r}.png`"
          width="70"
          height="70"
          class="mr-2"
        />
      </div>
      <div class="d-flex mt-4 bg-runes perks">
        <img
          v-for="(r, i) in page.perks.rows"
          :key="`perk-${r}-${i}`"
          :src="`/assets/rune-avis/perks/${r}.png`"
          width="50"
          height="50"
          class="perk mr-3"
          :class="`perk-${r}`"
        />
      </div>

      <div class="ctrl-btns">
        <button
          class="btn-slide mr-3 shadow"
          @click="savePage"
        >
          SAVE
        </button>
      </div>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../js/rest';
import EventBus from '../js/eventbus';

export default {
  name: 'Share',

  components: {},

  props: {},

  data: function() {
    return {
      ident: '',

      notfound: false,

      share: {},

      page: {
        title: '',
        champions: [],
        primary: {
          rows: [],
        },
        secondary: {
          rows: [],
        },
        perks: {
          rows: [],
        },
      },

      user: {},

      loggedin: false,
    };
  },

  created: function() {
    this.ident = this.$route.params.ident;
    Rest.getShare(this.ident)
      .then((res) => {
        if (res.body && res.body.share && res.body.page) {
          this.share = res.body.share;
          this.page = res.body.page;
          this.user = res.body.user;
        }
      })
      .catch((err) => {
        if (err.code === 404) {
          this.notfound = true;
        } else {
          console.error(err);
        }
      });

    Rest.getMe()
      .then(() => (this.loggedin = true))
      .catch(() => {});
  },

  methods: {
    savePage() {
      if (!this.loggedin) {
        this.$router.push({
          path: '/login',
          query: { createpage: this.ident },
        });
        return;
      }

      let page = JSON.parse(JSON.stringify(this.page));
      page.title += ` (cloned from ${this.user.displayname})`;

      Rest.createPage(page)
        .then((res) => {
          if (res.body) {
            this.$router.push({
              name: 'RunePage',
              params: { uid: res.body.uid },
            });
          }
        })
        .catch(console.error);
    },
  },
};
</script>

<style scoped>
/** @format */

.bg-runes {
  border-radius: 5px;
  width: fit-content;
  padding: 10px 10px 10px 20px;
}
</style>
