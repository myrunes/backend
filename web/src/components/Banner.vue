<!-- @format -->

<template>
  <div
    v-if="internalVisible === null ? visible : internalVisible"
    class="banner"
    :class="classObj"
    :style="{ width: width || '100%' }"
  >
    <div
      v-if="internalClosable === null ? closable : internalClosable"
      class="close"
      @click="hide(true)"
    ></div>
    <p
      v-if="text"
      class="m-0"
    >
      {{ text }}
    </p>
    <slot></slot>
  </div>
</template>

<script>
/** @format */

export default {
  name: 'Banner',

  props: {
    type: String,
    width: String,
    closable: Boolean,
    visible: Boolean,
  },

  data() {
    return {
      text: null,
      internalVisible: null,
      internalType: null,
      internalClosable: null,
    };
  },

  computed: {
    classObj: function() {
      let type =
        (this.internalType != null ? this.internalType : this.type) || 'info';
      let o = {};
      o[type] = true;
      return o;
    },
  },

  methods: {
    emitClosing(active) {
      this.$emit('closing', !!active);
    },

    show(type, text, time, closable) {
      this.internalVisible = true;
      this.internalType = type;
      this.text = text;

      if (this.closable != null) {
        this.internalClosable = closable;
      }

      if (time) {
        setTimeout(() => this.hide(), time);
      }
    },

    hide(active) {
      if (this.internalVisible) {
        this.internalVisible = false;
        this.emitClosing(active);
      }
    },
  },
};
</script>

<style scoped>
/** @format */

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
