<template>
  <div class="outer-container">
    <div class="searchbar mx-auto position-relative">
      <input type="text" class="tb tb-champ" autocomplete="off"/>
      <span class="tb tb-champ"></span>
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
