<template>
  <div>
    <div class="champ-header mb-4" v-if="champ">
      <img :src="`/assets/champ-avis/${champ}.png`" width="42" height="42"/>
      <h2>{{ champ.toUpperCase() }}</h2>
    </div>
    <div>
      <Page v-for="p in pages" 
        :key="p.uid" 
        :uid="p.uid"
        :title="p.title"
        :champs="p.champions.join(' ')"
        :primary="p.primary.tree"
        :secondary="p.secondary.tree"
        :prows="p.primary.rows.join(' ')"
        :srows="p.secondary.rows.join(' ')"
        :perks="p.perks.rows.join(' ')"
      />
    </div>
    <button class="btn-slide btn-new">+</button>
  </div>
</template>

<script>
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
      pages: [],
    }
  },

  methods: {

  },

  created: function() {
    this.champ = this.$route.params.champ;
    Rest.getPages().then((res) => {
      if (!res.body) return;
      this.pages = res.body.data
        .filter((p) => p.champions.includes(this.champ));
    }).catch(console.error);
  }
}

</script>


<style scoped>

.champ-header {
  display: flex;
}

.champ-header > img {
  border-radius: 50%;
  margin-right: 15px;
}

</style>
