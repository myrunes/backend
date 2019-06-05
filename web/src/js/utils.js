/** @format */

function parseTime(date) {
  function btf(inp) {
    if (inp < 10) return '0' + inp;
    return inp;
  }
  var date = date ? date : new Date();
  var y = date.getFullYear(),
    m = btf(date.getMonth() + 1),
    d = btf(date.getDate()),
    h = btf(date.getHours()),
    min = btf(date.getMinutes()),
    s = btf(date.getSeconds());
  return { y, m, d, h, min, s };
}

function getCookies() {
  var c = {};
  document.cookie
    .split(';')
    .map((v) => v.trim().split('='))
    .forEach((v) => (c[v[0]] = v[1]));
  return c;
}

function getCookieValue(name) {
  return getCookies()[name];
}

// ----------------------------------

export default {
  parseTime,
  getCookies,
  getCookieValue,
};
