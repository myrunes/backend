<!-- @format -->

<template>
  <div class="slider" @click="update">
    <span class="bar" :class="{ enabled: state }"></span>
    <span class="dot" :class="{ enabled: state }"></span>
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'Slider',

  data: function() {
    return {
      state: false,
    };
  },

  methods: {
    update() {
      this.state = !this.state;
      this.$emit('input', this.state);
    },

    set(val) {
      this.state = val;
    },
  },
};
</script>

<style scoped>
* {
  transition: all 0.25s ease-in-out;
}

.slider {
  display: flex;
}

.dot {
  position: relative;
  left: -2.8em;

  background-color: rgb(40, 102, 131);
  width: 1.3em;
  height: 1.3em;
  margin-top: auto;
  margin-bottom: auto;
  /* margin-right: 20px; */
  border-radius: 50%;
}

.bar {
  width: 3em;
  height: 1.7em;
  background-color: rgb(52, 67, 75);
  margin-top: auto;
  margin-bottom: auto;
  border-radius: 1.7em;
  border: solid rgba(3, 168, 244, 0) 2px;
}

.dot.enabled {
  background-color: rgba(3, 168, 244, 0.671);
  left: -1.6em;
}

.bar.enabled {
  border: solid #03a9f4 2px;
}

.slider:hover {
  cursor: pointer;
}
</style>
