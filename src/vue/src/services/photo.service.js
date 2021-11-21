//import config from 'config';
import { authHeader } from "../helpers";
var config = {};
config.apiUrl = window.location.origin;

export const photoService = {
  create,
  get,
  getSince,
  getAll
};

function getAll() {
  const requestOptions = {
    method: "GET",
    headers: authHeader()
  };

  return fetch(`${config.apiUrl}/api/photos`, requestOptions).then(
    handleResponse
  );
}

function get(photo_id) {
  const requestOptions = {
    method: "GET",
    headers: authHeader()
  };

  return fetch(`${config.apiUrl}/api/photo/${photo_id}`, requestOptions).then(
    handleResponse
  );
}

function getSince(timestamp) {
  var url = new URL(`/api/photos`, config.apiUrl);
  const params = { 
      timestamp: timestamp
  };
  const requestOptions = {
    method: "GET",
    headers: authHeader(),
  };
  url.search = new URLSearchParams(params).toString();
  return fetch(url, requestOptions).then(
    handleResponse
  );
}

function create(photoname) {
  const requestOptions = {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify({ name: photoname})
  };

  return fetch(`${config.apiUrl}/api/v1/new_photo`, requestOptions).then(
    handleResponse
  );
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
