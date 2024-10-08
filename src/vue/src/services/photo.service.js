//import config from 'config';
var config = {};
config.apiUrl = window.location.origin;

export const photoService = {
  get,
  getTicket,
  getAll,
  getOwn,
  putMetadata,
  put,
  del
};

function get(id) {
  const requestOptions = {
    method: "GET"
  };

  return fetch(`${config.apiUrl}/api/photo/${id}`, requestOptions).then(
    handleResponse
  );
}

function getTicket() {
  const requestOptions = {
    method: "GET"
  };

  return fetch(`${config.apiUrl}/api/new_photo`, requestOptions).then(
    handleResponse
  );
}

function getAll() {
  const requestOptions = {
    method: "GET"
  };
  return fetch(`${config.apiUrl}/api/events`, requestOptions).then(
    handleResponse
  );
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

function putMetadata(metadata) {
  var url = new URL(`/api/new_photo`, config.apiUrl);
  const params = {
    ticket_id: metadata.ticket_id
  };
  const requestOptions = {
    headers: {
      "Content-Type": "application/json"
    },
    method: "PUT",
    body: JSON.stringify({
      author_id: metadata.author_id,
      author: metadata.author,
      description: metadata.description
    })
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(handleResponse);
}

function put(uid, photo) {
  var url = new URL(`/api/photo/${photo.photo_id}`, config.apiUrl);
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
    if (typeof data.data !== "undefined") {
      return data.data;
    } else {
      return data;
    }
  });
}
