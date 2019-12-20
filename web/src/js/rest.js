/** @format */

import request from 'request';

const HOST =
  process.env.NODE_ENV === 'production'
    ? window.location.origin
    : 'http://localhost:8080';

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
    url: `${HOST}/api/users`,
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

function getChamps() {
  return _req({
    url: `${HOST}/api/resources/champions`,
    method: 'GET',
  });
}

function getRunes() {
  return _req({
    url: `${HOST}/api/resources/runes`,
    method: 'GET',
  });
}

function getPages(sortBy, champion, short, filter) {
  return _req({
    url: `${HOST}/api/pages`,
    method: 'GET',
    qs: {
      sortBy,
      champion,
      short,
      filter,
    },
  });
}

function getPage(uid) {
  return _req({
    url: `${HOST}/api/pages/${uid}`,
    method: 'GET',
  });
}

function updatePage(uid, page) {
  return _req({
    url: `${HOST}/api/pages/${uid}`,
    method: 'POST',
    json: page,
  });
}

function createPage(page) {
  return _req({
    url: `${HOST}/api/pages`,
    method: 'POST',
    json: page,
  });
}

function deletePage(uid) {
  return _req({
    url: `${HOST}/api/pages/${uid}`,
    method: 'DELETE',
  });
}

function updateUser(update) {
  return _req({
    url: `${HOST}/api/users/me`,
    method: 'POST',
    json: update,
  });
}

function deleteUser(currpassword) {
  return _req({
    url: `${HOST}/api/users/me`,
    method: 'DELETE',
    json: { currpassword },
  });
}

function getSessions() {
  return _req({
    url: `${HOST}/api/sessions`,
    method: 'GET',
  });
}

function deleteSession(sessionid) {
  return _req({
    url: `${HOST}/api/sessions/${sessionid}`,
    method: 'DELETE',
  });
}

function getFavorites() {
  return _req({
    url: `${HOST}/api/favorites`,
    method: 'GET',
  });
}

function setFavorites(favorites) {
  return _req({
    url: `${HOST}/api/favorites`,
    method: 'POST',
    json: { favorites },
  });
}

function getShare(ident) {
  return _req({
    url: `${HOST}/api/shares/${ident}`,
    method: 'GET',
  });
}

function createShare(share) {
  return _req({
    url: `${HOST}/api/shares`,
    method: 'POST',
    json: share,
  });
}

function updateShare(share) {
  return _req({
    url: `${HOST}/api/shares/${share.uid}`,
    method: 'POST',
    json: share,
  });
}

function deleteShare(share) {
  return _req({
    url: `${HOST}/api/shares/${share.uid}`,
    method: 'DELETE',
  });
}

function getVersion() {
  return _req({
    url: `${HOST}/api/version`,
    method: 'GET',
  });
}

function getAPIToken() {
  return _req({
    url: `${HOST}/api/apitoken`,
    method: 'GET',
  });
}

function generateAPIToken() {
  return _req({
    url: `${HOST}/api/apitoken`,
    method: 'POST',
  });
}

function deleteAPIToken() {
  return _req({
    url: `${HOST}/api/apitoken`,
    method: 'DELETE',
  });
}

function setPageOrder(pageorder, champion) {
  return _req({
    url: `${HOST}/api/users/me/pageorder`,
    method: 'POST',
    json: { pageorder },
    qs: { champion },
  });
}

function setMailAddress(mailaddress, reset) {
  if (reset === undefined || reset === null) {
    reset = false;
  }

  return _req({
    url: `${HOST}/api/users/me/mail`,
    method: 'POST',
    json: { mailaddress, reset },
  });
}

function confirmMail(token) {
  return _req({
    url: `${HOST}/api/users/me/mail/confirm`,
    method: 'POST',
    json: { token },
  });
}

function resetPassword(mailaddress) {
  return _req({
    url: `${HOST}/api/users/me/passwordreset`,
    method: 'POST',
    json: { mailaddress },
  });
}

function resetPasswordConfirm(token, new_password, page_names) {
  return _req({
    url: `${HOST}/api/users/me/passwordreset/confirm`,
    method: 'POST',
    json: { token, new_password, page_names },
  });
}

// ----------------------------

function _req(options) {
  return new Promise((resolve, rejects) => {
    options.withCredentials = true;
    request(options, (err, res, body) => {
      if (err) {
        rejects(err);
        return;
      }

      if (body && typeof body === 'string') body = JSON.parse(body);

      if (res.statusCode >= 400) {
        body._headers = res.headers;
        rejects(body);
        return;
      }
      resolve({ res, body });
    });
  });
}

export default {
  getMe,
  checkUsername,
  register,
  login,
  logout,
  getChamps,
  getRunes,
  getPages,
  getPage,
  updatePage,
  createPage,
  deletePage,
  updateUser,
  deleteUser,
  getSessions,
  deleteSession,
  getFavorites,
  setFavorites,
  getShare,
  createShare,
  updateShare,
  deleteShare,
  getVersion,
  getAPIToken,
  generateAPIToken,
  deleteAPIToken,
  setPageOrder,
  setMailAddress,
  confirmMail,
  resetPassword,
  resetPasswordConfirm,
};
