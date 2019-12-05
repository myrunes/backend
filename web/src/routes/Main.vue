<!-- @format -->

<template>
  <div>
    <div class="outer-container">
      <div class="searchbar mx-auto position-relative">
        <input type="text" class="tb tb-champ" autocomplete="off" @input="searchAndDisplay" />
        <span class="tb tb-champ"></span>
      </div>
    </div>
    <h3 class="mx-auto my-5 text-center" v-if="displayFavs">YOUR FAVORITES</h3>
    <div class="container mt-5 champs-container">
      <a
        v-for="c in displayedChamps"
        :key="c"
        @click="openChamp(c)"
        class="champ-tile"
        :class="{ 'no-pages': !pages[c] }"
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

export default {
  name: 'Main',

  props: {},

  components: {},

  data: function() {
    return {
      champs: [],
      displayedChamps: [],
      pages: {},
      favorites: [],
      displayFavs: false,
    };
  },

  methods: {
    searchAndDisplay(e) {
      if (!e.target || e.target.value === undefined) return;
      let val = e.target.value.toLowerCase();
      if (val === '') {
        this.displayedChamps = this.favorites || [];
        this.displayFavs = this.favorites && this.favorites.length != 0;
      } else {
        console.log(this.champs.filter((c) => c.includes(val)));
        this.displayedChamps = this.champs.filter((c) => c.includes(val));
        this.displayFavs = false;
      }
    },

    openChamp(champ) {
      this.$router.push({ name: 'Champ', params: { champ } });
    },
  },

  created: function() {
    Rest.getChamps()
      .then((r) => {
        if (r.body && r.body.data) {
          this.champs = r.body.data;
        }
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
};
</script>

<style scoped>
/** @format */

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
  margin-top: 20vh;
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
</style>
