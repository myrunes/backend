<!-- @format -->

<template>
  <div>
    <div class="outer-container">
      <div class="searchbar mx-auto position-relative">
        <img class="search-icon" src="/assets/search.svg" />
        <input
          type="text"
          class="tb tb-champ"
          autocomplete="off"
          :placeholder="searchBarOnFocus ? '' : 'Search for a champ or page'"
          @focus="searchBarOnFocus = true"
          @blur="searchBarOnFocus = false"
          @input="searchAndDisplay"
        />
        <span class="tb tb-champ"></span>
      </div>
    </div>

    <h3 v-if="displayFavs" class="mx-auto my-5 text-center">YOUR FAVORITES</h3>

    <div class="container mt-5 champs-container">
      <div
        v-if="(!displayedChamps || displayedChamps.length < 1) && !isSearch"
        class="favorites-hint"
      >
        <img src="/assets/fav.svg" />
        <span>
          <h4>Did you know?</h4>
          <p>You can favorize champions which then are displayed here.</p>
        </span>
      </div>

      <a
        v-for="c in displayedChamps"
        :key="c"
        class="champ-tile"
        :class="{ 'no-pages': !pages[c] }"
        @click="openChamp(c)"
      >
        <img :src="`/assets/champ-avis/${c}.png`" width="100" height="100" />
        <p>{{ pages[c] }}</p>
      </a>
    </div>

    <div v-if="searchedPages" class="pages-container mt-5">
      <Page
        v-for="p in searchedPages"
        class="page-tile"
        :key="p.uid"
        :uid="p.uid"
        :title="p.title"
        :champs="p.champions"
        :primary="p.primary.tree"
        :secondary="p.secondary.tree"
        :prows="p.primary.rows"
        :srows="p.secondary.rows"
        :perks="p.perks.rows"
        :displayDelete="false"
      />
    </div>

    <div class="ctrl-btns">
      <button
        class="btn-slide btn-new"
        @click="$router.push({ name: 'RunePage', params: { uid: 'new' } })"
      >+</button>
    </div>
  </div>
</template>

<script>
/** @format */

import EventBus from '../js/eventbus';
import Rest from '../js/rest';
import Page from '../components/Page';

const SHORTS = {
  mf: 'miss-fortune',
  ww: 'warwick',
  lb: 'leblanc',
  tf: 'twisted-fate',
  gp: 'gankplank',
};

export default {
  name: 'Main',

  components: {
    Page,
  },

  data: function() {
    return {
      champs: [],
      displayedChamps: [],
      pages: {},
      favorites: [],
      displayFavs: false,
      searchedPages: null,
      isSearch: false,
      searchBarOnFocus: false,
    };
  },

  created: function() {
    Rest.getChamps()
      .then((res) => {
        this.champs = res.body.data;
      })
      .catch(console.error);

    Rest.getFavorites()
      .then((r) => {
        if (r.body && r.body.data) {
          this.favorites = r.body.data;
          this.displayedChamps = r.body.data;
          this.displayFavs = this.favorites && this.favorites.length != 0;
        }
      })
      .catch(console.error);

    Rest.getPages(null, null, true)
      .then((r) => {
        if (r.body && r.body.data) {
          this.pages = r.body.data;
        }
      })
      .catch(console.error);
  },

  methods: {
    searchAndDisplay(e) {
      if (!e.target || e.target.value === undefined) return;

      let val = e.target.value.toLowerCase();

      if (val === '') {
        this.displayedChamps = this.favorites || [];
        this.displayFavs = this.favorites && this.favorites.length != 0;
        this.searchedPages = null;
        this.isSearch = false;
        return;
      }

      this.displayedChamps = this.champs
        .filter((c) => this.searchFilter(c, val))
        .map((c) => c.uid);
      this.displayFavs = false;
      this.isSearch = true;

      setTimeout(
        (v, t) => {
          if (v === t.value) {
            this.searchAndDisplayPages(v);
          }
        },
        400,
        val,
        e.target
      );
    },

    searchAndDisplayPages(filter) {
      Rest.getPages('created', null, false, filter)
        .then((res) => {
          this.searchedPages = res.body.data;
        })
        .catch(console.error);
    },

    searchFilter(c, val) {
      const name = c.name.toLowerCase();
      const nameStripped = c.uid.replace('-', ' ');
      const nameConcat = c.uid.replace('-', '');

      val = val.toLowerCase();

      return (
        name.includes(val) ||
        nameStripped.includes(val) ||
        nameConcat.includes(val) ||
        SHORTS[val] === c.uid
      );
    },

    openChamp(champ) {
      this.$router.push({ name: 'Champ', params: { champ } });
    },
  },
};
</script>

<style scoped>
/** @format */

@keyframes favorites-hint-in {
  0% {
    opacity: 0;
  }
  33% {
    opacity: 0;
  }
  100% {
    opacity: 0.75;
  }
}

.outer-container {
  display: flex;
  width: 100%;
}

.tb-champ {
  width: 50vw;
  font-size: 30px;
  font-family: 'Montserrat', sans-serif;
}

span.tb-champ {
  width: 50vw;
}

.searchbar {
  position: relative;
  margin-top: 20vh;
}

.search-icon {
  position: absolute;
  height: 30px;
  width: 30px;
  top: 10px;
  left: -40px;
  opacity: 0.8;
}

.pages-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  max-width: 1000px;
  margin: 0 auto;
}

.page-tile {
  max-width: 370px;
  width: 370px;
  margin-right: 10px;
}

.champs-container {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
}

.champ-tile {
  transition: all 0.25s ease;
  position: relative;
}

.champ-tile > img {
  border: solid 2px rgba(3, 168, 244, 0);
}

.champ-tile > p {
  position: absolute;
  top: 5px;
  left: 5px;
  font-size: 12px;
  padding: 0px 6px;
  background-color: rgba(3, 168, 244, 0.75);
  opacity: 0;
  transition: all 0.25s ease;
}

.champ-tile:hover > img {
  border: solid 2px #03a9f4;
  transition: all 0.25s ease-in-out;
}

.champ-tile:hover > p {
  opacity: 1;
}

.no-pages {
  opacity: 0.6;
  filter: grayscale(0.9);
}

a:hover {
  cursor: pointer;
  opacity: 1;
  filter: none;
}

.favorites-hint {
  display: flex;
  max-width: 400px;
  opacity: 0.75;

  animation: favorites-hint-in 3s ease;
}

.favorites-hint > img {
  width: 85px;
  height: 85px;
  border: dashed 3px white;
  padding: 15px;
  margin-right: 30px;
}
</style>
