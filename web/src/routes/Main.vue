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
          placeholder="Search for a champion"
          @input="searchAndDisplay"
        />
        <span class="tb tb-champ"></span>
      </div>
    </div>

    <h3 v-if="displayFavs" class="mx-auto my-5 text-center">YOUR FAVORITES</h3>

    <div class="container mt-5 champs-container">
      <div v-if="!displayedChamps || displayedChamps.length < 1" class="favorites-hint">
        <img src="/assets/fav.svg" />
        <span>
          <h4>Did you know?</h4>
          <p>You can favorite champions which then are displayed here.</p>
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
  </div>
</template>

<script>
/** @format */

import EventBus from '../js/eventbus';
import Rest from '../js/rest';

const SHORTS = {
  mf: 'miss-fortune',
  ww: 'warwick',
  lb: 'leblanc',
  tf: 'twisted-fate',
  gp: 'gankplank',
};

export default {
  name: 'Main',

  components: {},

  props: {},

  data: function() {
    return {
      champs: [],
      displayedChamps: [],
      pages: {},
      favorites: [],
      displayFavs: false,
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
      } else {
        this.displayedChamps = this.champs
          .filter((c) => this.searchFilter(c, val))
          .map((c) => c.uid);
        this.displayFavs = false;
      }
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

.searchbar {
  position: relative;
  margin-top: 20vh;
}

.search-icon {
  position: absolute;
  height: 30px;
  width: 30px;
  top: 10px;
  opacity: 0.8;
}

span.tb-champ {
  width: 50vw;
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
