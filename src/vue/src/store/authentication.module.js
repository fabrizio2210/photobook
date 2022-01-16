import { uidService } from "../services";
import Vue from "vue";

const user = JSON.parse(localStorage.getItem("user"));
const initialState = user ? { status: {}, user } : { status: {}, user: null };

export const authentication = {
  namespaced: true,
  state: initialState,
  actions: {
    get_uid({ commit }) {
      uidService.getUid().then(uid => {
        commit("get_uid", uid["uid"]);
      });
    },
    set_name({ commit }, author) {
      commit("set_name", author);
    }
  },
  mutations: {
    get_uid(state, uid) {
      if (!state.user) {
        state.user = { uid: uid };
      } else {
        Vue.set(state.user, "uid", uid);
      }
      localStorage.setItem("user", JSON.stringify(state.user));
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
