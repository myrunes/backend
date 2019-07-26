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

function setWindowListener(event, cb) {
  if (typeof event === 'object') {
    event.forEach((e) => {
      setWindowListener(e, cb);
    });
  } else {
    window.addEventListener(event, cb);
  }
}

function removeWindowListener(event, cb) {
  if (typeof event === 'object') {
    event.forEach((e) => {
      removeWindowListener(e, cb);
    });
  } else {
    window.removeEventListener(event, cb);
  }
}

function copyToClipboard(text) {
  return new Promise((resolve, reject) => {
    var id = 'hidden-clipboard-area';
    var existsTextarea = document.getElementById(id);
    if (!existsTextarea) {
      var textarea = document.createElement('textarea');
      textarea.id = id;
      document.querySelector('body').appendChild(textarea);
      existsTextarea = document.getElementById(id);
    }
    existsTextarea.value = text;
    existsTextarea.select();
    var status = document.execCommand('copy');
    if (!status) {
      reject('Could not copy shortlink to clipboard.');
    } else {
      resolve();
    }
  });
}

// ----------------------------------

export default {
  parseTime,
  getCookies,
  getCookieValue,
  setWindowListener,
  removeWindowListener,
  copyToClipboard,
};
