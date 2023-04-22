//import config from 'config';
var config = {};
config.apiUrl = window.location.origin;

export const photoService = {
  create,
  get,
  getSince,
  getOwn,
  getAll,
  put,
  del
};

function getAll() {
  const requestOptions = {
    method: "GET"
  };

  return fetch(`${config.apiUrl}/api/events`, requestOptions).then(
    handleResponse
  );
}

function get(id) {
  const requestOptions = {
    method: "GET"
  };

  return fetch(`${config.apiUrl}/api/photo/${id}`, requestOptions).then(
    handleResponse
  );
}

function getSince(timestamp) {
  var url = new URL(`/api/events`, config.apiUrl);
  const params = {
    timestamp: timestamp
  };
  const requestOptions = {
    method: "GET"
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function getOwn(uid) {
  var url = new URL(`/api/events`, config.apiUrl);
  const params = {
    author_id: uid
  };
  const requestOptions = {
    method: "GET"
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function create(photoname) {
  const requestOptions = {
    method: "POST",
    body: JSON.stringify({ name: photoname })
  };

  return fetch(`${config.apiUrl}/api/new_photo`, requestOptions).then(
    handleResponse
  );
}

function put(uid, photo) {
  var url = new URL(`/api/photo/${photo.id}`, config.apiUrl);
  const params = {
    author_id: uid
  };
  const requestOptions = {
    headers: {
      "Content-Type": "application/json"
    },
    method: "PUT",
    body: JSON.stringify({
      author: photo.author,
      description: photo.description
    })
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function del(uid, id) {
  var url = new URL(`/api/photo/${id}`, config.apiUrl);
  const params = {
    author_id: uid
  };
  const requestOptions = {
    method: "DELETE"
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function handleResponse(response) {
  return response.text().then(text => {
    const data = text && JSON.parse(text);
    if (!response.ok) {
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    return data;
  });
}
