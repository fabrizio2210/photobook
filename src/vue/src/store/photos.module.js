import { photoService } from "../services";
import Vue from "vue";
import router from "../router/";

export const photos = {
  namespaced: true,
  state: {
    all: { photos_list: [] },
    my: { photos_list: [] },
    last_timestamp: 0,
    status: {}
  },
  actions: {
    creating({ dispatch, commit }, { photoname }) {
      commit("creatingRequest", {});
      photoService.create(photoname).then(
        photo => {
          commit("creatingSuccess", photo);
          router.push({
            name: "PhotoSettings",
            params: { photo_id: photo.id }
          });
        },
        error => {
          commit("creatingFailure", error);
          dispatch("alert/error", error, { root: true });
        }
      );
    },
    prepareEdit({ commit }, { id }) {
      commit("prepareEdit", id);
    },
    edit({ commit }, { uid, photo }) {
      photoService.put(uid, photo).then(
        photo => commit("editSuccess", photo["photo"]),
        error => commit("editFailure", error)
      );
    },
    unedit({ commit }, { id }) {
      commit("unedit", id);
    },
    getAll({ commit }) {
      commit("getAllRequest");

      photoService.getAll().then(
        photos => commit("getAllSuccess", photos["photos"]),
        error => commit("getAllFailure", error)
      );
    },
    getSince({ commit }, { last_timestamp }) {
      commit("getSinceRequest");

      photoService.getSince(last_timestamp).then(
        photos => commit("getSinceSuccess", photos["photos"]),
        error => commit("getSinceFailure", error)
      );
    },
    get({ commit }, { id }) {
      commit("getRequest");

      photoService.get(id).then(
        photo => commit("getSuccess", photo["photo"]),
        error => commit("getFailure", { error, id })
      );
    },
    getOwn({ commit }, { uid }) {
      commit("getOwnRequest");

      photoService.getOwn(uid).then(
        photos => commit("getOwnSuccess", photos["photos"]),
        error => commit("getOwnFailure", error)
      );
    },
    del({ commit }, { uid, id }) {
      commit("deleteRequest");
      photoService.del(uid, id).then(
        photo => commit("deleteSuccess", photo["photo"]),
        error => commit("deleteFailure", error)
      );
    }
  },
  mutations: {
    creatingRequest(state) {
      state.status = { creating: true };
    },
    creatingSuccess(state, photo) {
      state.status = { created: true };
      state.all.photos_list.push(photo);
    },
    creatingFailure(state) {
      state.status = {};
      state.user = null;
    },
    getAllRequest(state) {
      state.all = { loading: true };
    },
    getRequest(state) {
      state.all.loading = true;
    },
    getSinceRequest(state) {
      state.all.loading = true;
    },
    deleteRequest(state) {
      state.all.loading = true;
    },
    getAllSuccess(state, photos) {
      const photos_list = photos.slice().reverse();
      Vue.set(state, "last_timestamp", photos_list[0].timestamp);
      state.all = { photos_list };
    },
    getSinceSuccess(state, photos) {
      state.all.photos_list = photos
        .slice()
        .reverse()
        .concat(state.all.photos_list);
      if (state.all.photos_list.length > 0) {
        state.last_timestamp = state.all.photos_list[0].timestamp;
      }
      Vue.delete(state.all, "loading");
    },
    getOwnRequest(state) {
      state.all.loading = true;
    },
    getSuccess(state, photo) {
      if (state.all.photos_list.length > 0) {
        for (var i = 0; i < state.all.photos_list.length; i++) {
          if (state.all.photos_list[i].id == photo.id) {
            state.all.photos_list[i] = photo;
          }
        }
      } else {
        var photos_list = [photo];
        state.all = { photos_list };
      }
      Vue.delete(state.all, "loading");
    },
    getOwnSuccess(state, photos) {
      const photos_list = photos.slice().reverse();
      state.my = { photos_list };
      Vue.delete(state.all, "loading");
    },
    deleteSuccess(state, photo) {
      if (typeof state.all.photos_list !== "undefined") {
        for (var i = 0; i < state.all.photos_list.length; i++) {
          if (state.all.photos_list[i].id == photo.id) {
            state.all.photos_list.splice(i, 1);
          }
        }
      }
      if (typeof state.my.photos_list !== "undefined") {
        for (i = 0; i < state.my.photos_list.length; i++) {
          if (state.my.photos_list[i].id == photo.id) {
            state.my.photos_list.splice(i, 1);
          }
        }
      }
      Vue.delete(state.all, "loading");
    },
    prepareEdit(state, id) {
      for (var i = 0; i < state.my.photos_list.length; i++) {
        if (state.my.photos_list[i].id == id) {
          Vue.set(state.my.photos_list[i], "edit", true);
          Vue.set(
            state.my.photos_list[i],
            "old_author",
            state.my.photos_list[i].author
          );
          Vue.set(
            state.my.photos_list[i],
            "old_description",
            state.my.photos_list[i].description
          );
        }
      }
    },
    editSuccess(state, photo) {
      for (var i = 0; i < state.my.photos_list.length; i++) {
        if (state.my.photos_list[i].id == photo.id) {
          Vue.set(state.my.photos_list, i, photo);
        }
      }
      if (
        typeof state.all.photos_list !== "undefined" &&
        state.all.photos_list.length > 0
      ) {
        for (i = 0; i < state.all.photos_list.length; i++) {
          if (state.all.photos_list[i].id == photo.id) {
            state.all.photos_list[i] = photo;
          }
        }
      } else {
        var photos_list = [photo];
        state.all = { photos_list };
      }
    },
    unedit(state, id) {
      for (var i = 0; i < state.my.photos_list.length; i++) {
        if (state.my.photos_list[i].id == id) {
          Vue.set(state.my.photos_list[i], "edit", false);
          Vue.set(
            state.my.photos_list[i],
            "author",
            state.my.photos_list[i].old_author
          );
          Vue.set(
            state.my.photos_list[i],
            "description",
            state.my.photos_list[i].old_description
          );
        }
      }
    },
    getAllFailure(state, error) {
      state.all = { error };
    },
    getSinceFailure(state, error) {
      state.all = { error };
    },
    getOwnFailure(state, error) {
      state.all = { error };
    },
    getFailure(state, error) {
      if (error.error == "Item not found.") {
        for (var i = 0; i < state.all.photos_list.length; i++) {
          if (state.all.photos_list[i].id == error.id) {
            state.all.photos_list.splice(i, 1);
          }
        }
        Vue.delete(state.all, "loading");
      } else {
        state.all = { error };
      }
    },
    deleteFailure(state, error) {
      state.all = { error };
    }
  }
};
