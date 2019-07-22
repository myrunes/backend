<!-- @format -->

<template>
  <div>
    <div class="d-flex flex-wrap">
      <div v-for="e in elements" :key="`element-${e}`" class="element mb-2">
        <p>{{ e }}</p>
        <a class="xbtn" @click="removeElement(e, true)">X</a>
      </div>
      <input type="text" class="tags-tb mb-2" ref="tbInput" @input="tbInput" />
    </div>
    <div class="d-flex flex-wrap mt-3">
      <a
        v-for="s in suggestions"
        :key="`element-${s}`"
        class="suggestion mb-2"
        @click="append(s, true)"
      >
        {{ s }}
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TagsInput',

  props: {
    tags: {
      type: Array,
      default: () => {
        return [];
      },
    },
  },

  data: function() {
    return {
      elements: [],
      suggestions: [],
    };
  },

  methods: {
    append(e, emit) {
      this.elements.push(e);
      this.$refs.tbInput.value = '';
      this.suggestions = [];
      if (emit) this.emitAppend(e);
    },

    removeElement(e, emit) {
      let i = this.elements.findIndex((el) => el === e);
      if (i >= 0) {
        this.elements.splice(i, 1);
      }
      if (emit) this.emitRemove(e);
    },

    tbInput(e) {
      let val = e.target.value.toLowerCase();
      if (!val) {
        this.suggestions = [];
      } else {
        this.suggestions = this.tags.filter(
          (t) => t.includes(val) && !this.elements.find((e) => e.includes(val))
        );
      }
    },

    emitAppend(e) {
      this.$emit('input', e);
      this.$emit('change', this.elements);
    },

    emitRemove(e) {
      this.$emit('remove', e);
      this.$emit('change', this.elements);
    },
  },
};
</script>

<style scoped>
p {
  margin: 0px;
}

.element {
  display: flex;
  margin-right: 10px;
  background-color: #37474f;
  padding: 5px 10px;
  border-radius: 5px;
}

.suggestion {
  margin-right: 10px;
  background-color: #0277bd;
  padding: 5px 10px;
  border-radius: 5px;
}

.xbtn {
  padding: 3px 8px;
  background-color: #212121;
  font-size: 12px;
  border-radius: 5px;
  margin-left: 8px;
  font-family: 'Roboto', sans-serif;
  font-weight: 500;
}

.tags-tb {
  border: solid 2px #0277bd;
  background-color: #37474f;
  color: white;
  border-radius: 5px;
  padding: 0px 5px;
  width: 150px;
}
</style>
