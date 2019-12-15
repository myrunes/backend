<template>
  <div
    v-if="wrapper"
    class="wrapper"
  >
    <div
      class="bubble"
      :style="{ 'background-image': gradient }"
      :class="{ visible: bubble }"
    >
      <div class="my-auto">
        <slot></slot>
      </div>
      <a
        class="close-btn"
        @click="hide(true)"
      >
        <img src="/assets/close.svg" />
      </a>
    </div>
  </div>
</template>

<script>
/** @format */

export default {
  name: 'InfoBubble',

  props: {
    color: String,
  },

  data() {
    return {
      wrapper: false,
      bubble: false,
    };
  },

  computed: {
    gradient() {
      let color1 = '#33691E';
      let color2 = '#9E9D24';
      switch (this.color) {
        case 'red':
          color1 = '#b71c1c';
          color2 = '#f44336';
          break;
        case 'cyan':
          color1 = '#0D47A1';
          color2 = '#03A9F4';
          break;
        case 'orange':
          color1 = '#E65100';
          color2 = '#FFA000';
          break;
      }

      return `linear-gradient(-10deg, ${color1}, ${color2})`;
    },
  },

  methods: {
    show() {
      this.wrapper = true;
      setTimeout(() => {
        this.bubble = true;
      }, 10);
      this.$emit('shows');
    },

    hide(selfTriggered) {
      this.bubble = false;
      setTimeout(() => {
        this.wrapper = false;
      }, 700);
      this.$emit('hides', !!selfTriggered);
    },
  },
};
</script>

<style scoped>
/** @format */

.wrapper {
  display: flex;
  position: fixed;
  left: 0;
  right: 0;
  bottom: 70px;
}

.bubble {
  display: flex;
  width: fit-content;
  height: fit-content;
  padding: 6px 8px 6px 14px;
  border-radius: 1000px;
  box-shadow: 0px 0px 20px 0px rgba(0, 0, 0, 0.4);
  z-index: 5;
  margin: 0 auto;

  opacity: 0;
  transform: translateY(50px);
  transition: all 0.7s cubic-bezier(0, 0, 0, 1.26);
}

.bubble.visible {
  opacity: 1;
  transform: translateY(0px);
}

.bubble * {
  margin: 0;
}

.close-btn {
  margin: auto 0px auto 10px;
  opacity: 0.8;
  transition: all 0.25s ease-in-out;
}

.close-btn:hover {
  opacity: 1;
}

.close-btn > img {
  height: 25px;
  width: 25px;
}
</style>

