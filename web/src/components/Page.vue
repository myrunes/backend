<!-- @format -->

<template>
  <div class="page" @click="toPage(uid, $event)">
    <div>
      <div class="headline">
        <h3 class="mb-3 mr-3">{{ title }}</h3>
        <div class="champs mt-2">
          <img
            v-for="c in champs.split(' ')"
            :key="c"
            :src="`/assets/champ-avis/${c}.png`"
            width="20"
            height="20"
          />
        </div>
      </div>
      <div class="runes" :class="`tree-${primary}`">
        <img
          :src="`/assets/rune-avis/${primary}.png`"
          class="mr-3"
          width="50"
          height="50"
        />
        <img
          v-for="r in prows.split(' ')"
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
          v-for="r in srows.split(' ')"
          :key="r"
          :src="`/assets/rune-avis/${secondary}/${r}.png`"
          width="50"
          height="50"
        />
      </div>
      <div class="runes perks">
        <img
          v-for="(r, i) in perks.split(' ')"
          :key="`perk-${r}-${i}`"
          :src="`/assets/rune-avis/perks/${r}.png`"
          width="30"
          height="30"
          class="perk"
          :class="`perk-${r}`"
        />
      </div>
    </div>
    <div>
      <button class="btn-slide btn-delete" @click="deletePage">
        {{ suredel ? 'SURE?' : 'DELETE' }}
      </button>
    </div>
  </div>
</template>

<script>
import Rest from '../js/rest';

export default {
  name: 'Page',

  props: {
    uid: String,
    title: String,
    champs: String,
    primary: String,
    secondary: String,
    prows: String,
    srows: String,
    perks: String,
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
* {
  transition: all 0.25s ease-in-out;
}

.page {
  display: flex;
  justify-content: space-between;

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
}

.btn-delete::before {
  background-color: rgb(229, 57, 53);
}
</style>
