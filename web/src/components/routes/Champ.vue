<!-- @format -->

<template>
  <div>
    <div class="champ-header mb-4" v-if="champ">
      <img :src="`/assets/champ-avis/${champ}.png`" width="42" height="42" />
      <h2>{{ champ.toUpperCase() }}</h2>
    </div>
    <div>
      <Page
        v-for="p in pages"
        :key="p.uid"
        :uid="p.uid"
        :title="p.title"
        :champs="p.champions.join(' ')"
        :primary="p.primary.tree"
        :secondary="p.secondary.tree"
        :prows="p.primary.rows.join(' ')"
        :srows="p.secondary.rows.join(' ')"
        :perks="p.perks.rows.join(' ')"
        @delete="deleted"
      />
    </div>
    <div class="ctrl-btns">
      <button
        class="btn-slide btn-new favorite"
        :class="{ active: this.favorites.includes(this.champ) }"
        @click="toggleFav"
      ></button>
      <button
        class="btn-slide btn-new"
        @click="
          $router.push({
            name: 'RunePage',
            params: { uid: 'new' },
            query: { champ },
          })
        "
      ></button>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../../js/rest';
import Page from '../Page';

export default {
  name: 'Champ',

  components: {
    Page,
  },

  data: function() {
    return {
      champ: null,
      favorite: false,
      favorites: [],
      pages: [],
    };
  },

  methods: {
    reload() {
      Rest.getPages()
        .then((res) => {
          if (!res.body) return;
          this.pages = res.body.data.filter((p) =>
            p.champions.includes(this.champ)
          );

          Rest.getFavorites()
            .then((res) => {
              if (!res.body || !res.body.data) return;
              this.favorites = res.body.data;
            })
            .catch(console.error);
        })
        .catch(console.error);
    },

    deleted() {
      this.reload();
    },

    toggleFav() {
      let ind = this.favorites.indexOf(this.champ);
      if (ind > -1) {
        this.favorites.splice(ind, 1);
      } else {
        this.favorites.push(this.champ);
      }

      Rest.setFavorites(this.favorites).catch(console.error);
    },
  },

  created: function() {
    this.champ = this.$route.params.champ;
    this.reload();
  },
};
</script>

<style scoped>
/** @format */

.champ-header {
  display: flex;
}

.champ-header > img {
  border-radius: 50%;
  margin-right: 15px;
}

.favorite::after {
  background-image: url('/assets/fav.svg');
}

.active::after {
  background-image: url('/assets/fav-active.svg');
}
</style>
