import { uidService } from "../services";
import Vue from "vue";

const user = JSON.parse(localStorage.getItem("user"));
const initialState = user ? { status: {}, user } : { status: {}, user: null };

export const authentication = {
  namespaced: true,
  state: initialState,
  actions: {
    async get_uid({ commit }) {
      return uidService.getUid().then(uid => {
        commit("set_uid", uid["uid"]);
      });
    },
    set_uid({ commit }, id) {
      commit("set_uid", id["id"]);
    },
    set_name({ commit }, author) {
      commit("set_name", author);
    }
  },
  mutations: {
    set_uid(state, uid) {
      if (typeof uid != "undefined") {
        if (!state.user) {
          state.user = { uid: uid };
        } else {
          Vue.set(state.user, "uid", uid);
        }
        localStorage.setItem("user", JSON.stringify(state.user));
      }
    },
    set_name(state, author) {
      if (!state.user) {
        state.user = { name: author };
      } else {
        Vue.set(state.user, "name", author);
      }
      localStorage.setItem("user", JSON.stringify(state.user));
    }
  }
};
