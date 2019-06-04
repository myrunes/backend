<template>
  <div class="page">
    <div>
      <h3 class="mb-3">{{ title }}</h3>
      <div class="runes">
        <img :src="`/assets/rune-avis/${primary}.png`" 
          class="mr-3" width="50" height="50"
        />
        <img v-for="r in prows.split(' ')"
          :key="r"
          :src="`/assets/rune-avis/${primary}/${r}.png`"
          width="50" height="50"
        />
      </div>
      <div class="runes">
        <img :src="`/assets/rune-avis/${secondary}.png`" 
          class="mr-3" width="50" height="50"
        />
        <img v-for="r in srows.split(' ')"
          :key="r"
          :src="`/assets/rune-avis/${secondary}/${r}.png`"
          width="50" height="50"
        />
      </div>
      <div class="runes">
        <img v-for="r in perks.split(' ')"
          :key="r"
          :src="`/assets/rune-avis/perks/${r}.png`"
          width="30" height="30"
          class="perk"
          :class="`perk-${r}`"
        />
      </div>
      <div class="champs">
        <img v-for="c in champs.split(' ')"
          :key="c"
          :src="`/assets/champ-avis/${c}.png`"
          width="15" height="15"
        />
      </div>
    </div>
  </div>
</template>

<script>

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
}

</script>

<style scoped>

* {
  transition: all .25s ease-in-out;
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
  margin: 0px 0px 20px 0px;
}

.runes > img {
  margin-right: 10px;
}

.perk {
  padding: 5px;
  border-radius: 50%;
}

</style>
