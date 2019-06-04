<template>
  <div>
    
  </div>
</template>

<script>
import Rest from '../../js/rest';

export default {
  name: 'Edit',

  components: {
  },

  data: function() {
    return {
      uid: null,
      page: null,
    }
  },

  methods: {

  },

  created: function() {
    this.uid = this.$route.params.uid;
    Rest.getPage(this.uid).then((res) => {
      if (!res.body) return;
      this.page = res.body;
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
