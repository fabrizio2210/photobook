export function authHeader() {
  // return authorization header with jwt token
  let user = JSON.parse(localStorage.getItem("user"));

  if (user && user.uid) {
    return {
      Authorization: "JWT " + user.uid,
      "Content-Type": "application/json"
    };
  } else {
    return {};
  }
}
