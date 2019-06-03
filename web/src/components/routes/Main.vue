<template>
  <div>
    <div class="outer-container">
      <div class="searchbar mx-auto position-relative">
        <input type="text" class="tb tb-champ" autocomplete="off" @input="searchAndDisplay"/>
        <span class="tb tb-champ"></span>
      </div>
    </div>
    <div class="container mt-5">
      <a v-for="c in displayedChamps" :key="c">
        <img :src="`/assets/champ-avis/${c}.png`"/>
      </a>
    </div>
  </div>
</template>

<script>
import EventBus from '../../js/eventbus';
import Rest from '../../js/rest';

export default {
  name: 'Main',
  
  props: {
  },

  components: {
  },

  data: function() {
    return {
      champs: [],
      displayedChamps: [],
    }
  },

  methods: {
    searchAndDisplay(e) {
      console.log(e);
      if (!e.target || e.target.value === undefined) return;
      let val = e.target.value.toLowerCase();
      if (val === '') {
        this.displayedChamps = [];
      } else {
        this.displayedChamps = this.champs.filter((c) => c.includes(val));
      }
    }
  },

  created: function() {
    Rest.getChamps().then((r) => {
      if (r.body && r.body.data) {
        this.champs = r.body.data;
      }
    }).catch(console.error);
  }

}
</script>

<style scoped>

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

</style>
