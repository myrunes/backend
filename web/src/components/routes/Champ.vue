<!-- @format -->

<template>
  <div>
    <SearchBar v-if="search" class="searchbar" @close="search = false" @input="onSearchInput" />
    <InfoBubble ref="info" color="orange" @hides="onInfoClose">
      <p>
        Searching for a specific page? Press
        <b>CTRL + F</b>!
      </p>
    </InfoBubble>
    <div class="champ-header mb-4" v-if="champ" :style="{ 'padding-top': search ? '20px' : '0' }">
      <img :src="`/assets/champ-avis/${champ}.png`" width="42" height="42" />
      <h2>{{ champ.toUpperCase() }}</h2>
    </div>
    <div>
      <Page
        v-for="p in pagesVisible"
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
import Utils from '../../js/utils';
import Page from '../Page';
import SearchBar from '../SearchBar';
import InfoBubble from '../InfoBubble';

export default {
  name: 'Champ',

  components: {
    Page,
    InfoBubble,
    SearchBar,
  },

  data: function() {
    return {
      champ: null,
      favorite: false,
      favorites: [],
      pages: [],
      pagesVisible: [],
      search: false,
    };
  },

  methods: {
    reload() {
      Rest.getPages()
        .then((res) => {
          if (!res.body) return;
          this.pages = this.pagesVisible = res.body.data.filter((p) =>
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

    onSearchPress(event) {
      if (event.keyCode == 114 || (event.ctrlKey && event.keyCode == 70)) {
        this.search = true;
        event.preventDefault();
      }
    },

    onEscapePress(event) {
      if (event.keyCode == 27) {
        this.search = false;
        this.pagesVisible = this.pages;
        event.preventDefault();
      }
    },

    onSearchInput(e) {
      let txt = e.text.toLowerCase();
      if (txt.length === 0) {
        this.pagesVisible = this.pages;
        return;
      }
      this.pagesVisible = this.pages.filter((p) => {
        return p.title.toLowerCase().includes(txt);
      });
    },

    onInfoClose(selfTriggered) {
      if (selfTriggered) {
        window.localStorage['info-page-search'] = '1';
      }
    },
  },

  created: function() {
    this.champ = this.$route.params.champ;
    this.reload();

    Utils.setWindowListener('keydown', this.onSearchPress);
    Utils.setWindowListener('keydown', this.onEscapePress);
  },

  destroyed: function() {
    Utils.removeWindowListener('keydown', this.onSearchPress);
    Utils.removeWindowListener('keydown', this.onEscapePress);
  },

  mounted() {
    if (!window.localStorage['info-page-search']) {
      setTimeout(this.$refs.info.show, 3000);
    }
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
