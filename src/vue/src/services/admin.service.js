//import config from 'config';
var config = {};
config.apiUrl = window.location.origin;

export const adminService = {
  newPrint,
  getUpload,
  toggleUpload
};

function newPrint(uid) {
  var url = new URL(`/api/new_print`, config.apiUrl);
  const params = {
    author_id: uid
  };
  const requestOptions = {
    headers: {
      "Content-Type": "application/json"
    },
    method: "POST"
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function toggleUpload(uid) {
  var url = new URL(`/api/admin/toggle_upload`, config.apiUrl);
  const params = {
    author_id: uid
  };
  const requestOptions = {
    headers: {
      "Content-Type": "application/json"
    },
    method: "POST"
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function getUpload() {
  var url = new URL(`/api/admin/upload`, config.apiUrl);
  const params = {};
  const requestOptions = {
    headers: {
      "Content-Type": "application/json"
    },
    method: "GET"
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function handleResponse(response) {
  return response.text().then(text => {
    var data = text;
    try {
      data = JSON.parse(text);
    } catch {
      data = text;
    }
    if (!response.ok) {
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    if (typeof data.data !== "undefined") {
      return data.data;
    } else {
      return data;
    }
  });
}
