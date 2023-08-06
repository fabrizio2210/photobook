//import config from 'config';
var config = {};
config.apiUrl = window.location.origin;

export const adminService = {
  newPrint
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

function handleResponse(response) {
  return response.text().then(text => {
    const data = text && JSON.parse(text);
    if (!response.ok) {
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    if (typeof data.data !== 'undefined'){
      return data.data;
    } else {
      return data;
    }
  });
}
