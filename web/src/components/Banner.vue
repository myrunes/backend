<!-- @format -->

<template>
  <div class="banner" :class="classObj" :style="{ width: width || '100%' }">
    <div class="close" v-if="closable" @click="onclose"></div>
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'Banner',

  props: {
    type: String,
    width: String,
    closable: Boolean,
  },

  methods: {
    onclose() {
      this.$emit('closing');
    },
  },

  computed: {
    classObj: function() {
      let type = this.type || 'info';
      let o = {};
      o[type] = true;
      return o;
    },
  },
};
</script>

<style scoped>
.banner {
  padding: 10px 15px;
  border-width: 2px;
  border-style: solid;
}

.info {
  border-color: #00bcd4;
  background-color: #00bcd433;
}

.error {
  border-color: #f44336;
  background-color: #f4433633;
}

.success {
  border-color: #8bc34a;
  background-color: #8bc34a33;
}

.warning {
  border-color: #ffeb3b;
  background-color: #ffeb3b33;
}

.close {
  width: 1em;
  height: 1em;
  background-image: url('/assets/close.svg');
  background-size: 100%;
  cursor: pointer;
}
</style>
