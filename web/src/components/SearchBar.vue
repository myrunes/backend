<template>
  <div class="searchbar">
    <input ref="inpt" type="text" @input="onInput" />
    <a @click="onClose">
      <img src="/assets/close-black.svg" />
    </a>
  </div>
</template>

<script>
/** @format */

export default {
  name: 'SearchBar',

  methods: {
    onClose() {
      this.$emit('close');
    },

    onInput(event) {
      console.log(this.$refs);

      this.$emit('input', {
        text: event.target.value,
        event: event,
      });
    },
  },

  mounted() {
    this.$refs.inpt.focus();
  },
};
</script>

<style scoped>
/** @format */
.searchbar {
  display: flex;
  height: 50px;
  background-color: rgba(255, 255, 255, 0.8);
  border-radius: 5px;
}

.searchbar > input {
  width: 100%;
  margin: 12px 20px;
  background-color: transparent;
  border: none;
  border-bottom: 2px solid rgb(63, 63, 63);
}

.searchbar > a {
  margin: auto 20px auto 0px;
}

.searchbar > a > img {
  height: 30px;
  width: 30px;
}
</style>
