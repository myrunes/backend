/** @format */

import request from 'request';

const HOST =
  process.env.NODE_ENV === 'production' ? '' : 'http://localhost:8080';

function getMe() {
  return _req({
    url: `${HOST}/api/users/me`,
    method: 'GET',
    withCredentials: true,
  });
}

function checkUsername(uname) {
  return _req({
    url: `${HOST}/api/users/${uname}`,
    method: 'GET',
  });
}

function register(username, password, remember) {
  return _req({
    url: `${HOST}/api/users/me`,
    method: 'POST',
    json: {
      username,
      password,
      remember,
    },
  });
}

function login(username, password, remember) {
  return _req({
    url: `${HOST}/api/login`,
    method: 'POST',
    json: {
      username,
      password,
      remember,
    },
  });
}

function logout() {
  return _req({
    url: `${HOST}/api/logout`,
    method: 'POST',
  });
}

function _req(options) {
  return new Promise((resolve, rejects) => {
    options.withCredentials = true;
    request(options, (err, res, body) => {
      if (err) {
        rejects(err);
        return;
      }
      if (res.statusCode >= 400) {
        body = JSON.parse(body);
        rejects(body);
        return;
      }
      resolve(res, body);
    });
  });
}

export default {
  getMe,
  checkUsername,
  register,
  login,
  logout,
};
