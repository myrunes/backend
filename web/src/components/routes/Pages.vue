<template>
  <div>
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
        @delete="deleted"
      />
    </div>
    <button 
      class="btn-slide btn-new"
      @click="$router.push({ name: 'RunePage', params: { uid: 'new' } })"
    >
      +
    </button>
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
      pages: [],
    }
  },

  methods: {
    reload() {
      Rest.getPages().then((res) => {
        if (!res.body) return;
        this.pages = res.body.data;
      }).catch(console.error);
    },

    deleted() {
      this.reload();
    }
  },

  created: function() {
    this.reload();
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
