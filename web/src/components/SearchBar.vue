<template>
  <div class="searchbar">
    <img class="search-icon" src="/assets/search.svg" />
    <div>
      <input
        type="text"
        ref="inpt"
        class="tb tb-bar text-left"
        autocomplete="off"
        :placeholder="placeholder"
        @input="onInput"
      />
      <span class="tb tb-bar"></span>
    </div>
  </div>
</template>

<script>
/** @format */

export default {
  name: 'SearchBar',

  props: {
    placeholder: String,
  },

  mounted() {
    this.$refs.inpt.focus();
  },

  methods: {
    onInput(event) {
      this.$emit('input', {
        text: event.target.value,
        event: event,
      });
    },

    focus() {
      this.$refs.inpt.focus();
    },
  },
};
</script>

<style scoped>
/** @format */
.searchbar {
  position: relative;
  box-sizing: border-box;
  display: flex;
  width: 100%;
  /* height: 50px; */
}

.searchbar > div {
  position: relative;
  width: 100%;
  height: fit-content;
  margin: auto 0 auto 10px;
}

.tb-bar {
  width: 100%;
  /* font-size: 30px; */
  font-family: 'Montserrat', sans-serif;
}

span.tb-bar {
  width: 100%;
}

/* .searchbar > a {
  margin: auto 20px auto 0px;
} */

/* .searchbar > a > img {
  height: 30px;
  width: 30px;
} */
</style>
