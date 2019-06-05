<template>
  <footer>
    <p>
      © 2019 MYRUNES
      &nbsp;&nbsp;|&nbsp;&nbsp;
      <a href="/imprint.html">IMPRINT</a>
      &nbsp;&nbsp;|&nbsp;&nbsp;
      <a href="/">COOKIE USAGE</a>
      &nbsp;&nbsp;|&nbsp;&nbsp;
      <a href="https://github.com/zekroTJA/myrunes">GITHUB REPOSITORY</a>
      &nbsp;&nbsp;|&nbsp;&nbsp;
      CONTACT: contact[at]zekro.de
    </p>
  </footer>
</template>

<script>

export default {
  name: 'Footer',

  props: {
  },

  methods: {
  }
}
</script>

<style scoped>

p, a {
  margin: 0px;
  font-size: 14px;
  color: rgb(202, 202, 202) !important;
}

a {
  text-decoration: underline;
}

footer {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  padding: 5px 12px;
  background-color: #37474F;
  text-align: center;
  /* height: 60px; Set the fixed height of the footer here */
}

</style>
