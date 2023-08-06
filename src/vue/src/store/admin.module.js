import { adminService } from "../services";

export const admin = {
  namespaced: true,
  state: {
    status: {}
  },
  actions: {
    askPrint({ commit }, { uid }) {
      console.log(uid);
      commit("askPrintLoading");
      adminService.newPrint(uid).then(
        print => commit("askPrintSuccess", print),
        error => commit("askPrintFailure", error)
      );
    }
  },
  mutations: {
    askPrintLoading(state) {
      state.loading = true;
    },
    askPrintSuccess(state, print) {
      state.status = { print };
    },
    askPrintFailure(state, error) {
      state.error = { error };
    }
  }
};
