<template>
  <div 
    class="banner" 
    :class="classObj"
    :style="{ width: width || '100%' }"
  >
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'Banner',

  props: {
    type: String,
    width: String,
  },

  computed: {
    classObj: function() {
      let type = this.type || 'info';
      let o = {};
      o[type] = true;
      return o;
    }
  }
}
</script>

<style>

.banner {
  padding: 10px 15px;
  border-width: 2px;
  border-style: solid;
}

.info {
  border-color: #00BCD4;
  background-color: #00BCD433;
}

.error {
  border-color: #f44336;
  background-color: #f4433633;
}

.success {
  border-color: #8BC34A;
  background-color: #8BC34A33;
}

.warning {
  border-color: #FFEB3B;
  background-color: #FFEB3B33;
}


</style>
