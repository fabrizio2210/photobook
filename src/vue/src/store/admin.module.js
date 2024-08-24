import { adminService } from "../services";
import Vue from "vue";

export const admin = {
  namespaced: true,
  state: {
    print: {},
    upload: true
  },
  actions: {
    askPrint({ commit }, { uid }) {
      commit("loading");
      adminService.newPrint(uid).then(
        print => commit("setPrint", print),
        error => commit("failure", error)
      );
    },
    toggleUpload({ commit }, { uid }) {
      commit("loading");
      adminService.toggleUpload(uid).then(
        upload => commit("setUpload", upload),
        error => commit("failure", error)
      );
    },
    getUpload({ commit }) {
      commit("loading");
      adminService.getUpload().then(
        upload => commit("setUpload", upload),
        error => commit("failure", error)
      );
    }
  },
  mutations: {
    loading(state) {
      Vue.set(state, "loading", true);
      Vue.delete(state, "error");
    },
    failure(state, error) {
      Vue.delete(state, "loading");
      state.error = { error };
    },
    setPrint(state, print) {
      Vue.delete(state, "loading");
      state.print = { print };
    },
    setUpload(state, upload) {
      Vue.delete(state, "loading");
      Vue.set(state, "upload", upload["upload_status"]);
    }
  }
};
