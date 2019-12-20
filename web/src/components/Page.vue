<!-- @format -->

<template>
  <div class="page" @click="toPage(uid, $event)">
    <div class="headline">
      <h3 class="mb-3 mr-3">{{ title }}</h3>
      <div class="champs mt-1 mr-5">
        <img
          v-for="c in champs"
          :key="c"
          :src="`/assets/champ-avis/${c}.png`"
          width="25"
          height="25"
        />
      </div>
    </div>
    <div v-if="displayDelete">
      <button class="btn-slide btn-delete" @click="deletePage">
        {{ suredel ? 'SURE?' : 'DELETE' }}
      </button>
    </div>
    <div class="runes-container">
      <div class="runes" :class="`tree-${primary}`">
        <img
          :src="`/assets/rune-avis/${primary}.png`"
          class="mr-3"
          width="50"
          height="50"
        />
        <img
          v-for="r in prows"
          :key="r"
          :src="`/assets/rune-avis/${primary}/${r}.png`"
          width="50"
          height="50"
        />
      </div>
      <div class="runes" :class="`tree-${secondary}`">
        <img
          :src="`/assets/rune-avis/${secondary}.png`"
          class="mr-3"
          width="50"
          height="50"
        />
        <img
          v-for="r in srows"
          :key="r"
          :src="`/assets/rune-avis/${secondary}/${r}.png`"
          width="50"
          height="50"
        />
      </div>
      <div class="runes perks">
        <img
          v-for="(r, i) in perks"
          :key="`perk-${r}-${i}`"
          :src="`/assets/rune-avis/perks/${r}.png`"
          width="30"
          height="30"
          class="perk"
          :class="`perk-${r}`"
        />
      </div>
    </div>
  </div>
</template>

<script>
/** @format */

import Rest from '../js/rest';

export default {
  name: 'Page',

  props: {
    uid: String,
    title: String,
    champs: Array,
    primary: String,
    secondary: String,
    prows: Array,
    srows: Array,
    perks: Array,
    displayDelete: {
      type: Boolean,
      default: true,
    },
  },

  data: function() {
    return {
      suredel: false,
    };
  },

  methods: {
    toPage(uid, e) {
      if (e.target.localName === 'button') return;
      this.$router.push({ name: 'RunePage', params: { uid } });
    },

    deletePage() {
      if (this.suredel) {
        Rest.deletePage(this.uid)
          .then(() => {
            this.$emit('delete');
          })
          .catch(console.error);
      } else {
        this.suredel = true;
        setTimeout(() => {
          this.suredel = false;
        }, 2500);
      }
    },
  },
};
</script>

<style scoped>
/** @format */

* {
  transition: all 0.25s ease-in-out;
}

.page {
  display: grid;
  grid-template-columns: auto 100px;
  width: 100%;
  max-width: 100%;

  background-color: rgba(55, 71, 79, 0.726);
  margin-bottom: 10px;
  border-radius: 5px;
  padding: 15px 20px;
}

.page:hover {
  cursor: pointer;
  background-color: rgba(55, 71, 79, 1);
}

.champs {
  display: flex;
}

.champs > img {
  border-radius: 50%;
  margin-right: 5px;
}

.pathes {
  display: flex;
}

.pathes > img {
  border-radius: 50%;
  margin-right: 10px;
}

.runes-container {
  width: fit-content;
  grid-row-start: 2;
}

.runes {
  display: flex;
  margin: 0px 0px 10px 0px;
  padding: 10px;
  border-radius: 5px;
}

.runes > img {
  margin-right: 10px;
}

.headline {
  display: flex;
  white-space: nowrap;
  overflow: hidden;
}

.headline > h3 {
  overflow: hidden;
  text-overflow: ellipsis;
  display: block;
}

.btn-delete::before {
  background-color: rgb(229, 57, 53);
}
</style>
