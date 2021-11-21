import { userService } from "../services";
import router from "../router";
import Vue from "vue";

const user = JSON.parse(localStorage.getItem("user"));
const initialState = user
  ? { status: { loggedIn: true }, user }
  : { status: {}, user: null };

export const authentication = {
  namespaced: true,
  state: initialState,
  actions: {
    login({ dispatch, commit }, { username, password }) {
      commit("loginRequest", username );

      userService.login(username, password).then(
        user => {
          commit("loginSuccess", user);
          router.push("/");
        },
        error => {
          commit("loginFailure", error);
          dispatch("alert/error", error, { root: true });
        }
      );
    },
    token_refresh({ commit }, refresh_token) {
      userService.refresh(refresh_token).then(
        access_token => {
          commit("refresh", access_token['access_token']);
        }
      );
    },
    logout({ commit }) {
      userService.logout();
      commit("logout");
    }
  },
  mutations: {
    loginRequest(state, username) {
      state.status = { loggingIn: true };
      state.user = null;
      //TODO username should be get by the token
      state.username = username;
    },
    loginSuccess(state, user) {
      state.status = { loggedIn: true };
      state.user = user;
    },
    loginFailure(state) {
      state.status = {};
      state.user = null;
    },
    refresh(state, access_token) {
      Vue.set(state.user, 'access_token', access_token);
    },
    logout(state) {
      state.status = {};
      state.user = null;
    }
  }
};
