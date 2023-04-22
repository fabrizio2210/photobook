import { photoService } from "../services";
import Vue from "vue";

function removePhotoFromList(list, photo) {
  var found = false;
  for (var index in list) {
    if (list[index].photo_id == photo.photo_id) {
      found = true;
      break;
    }
  }
  if (found) {
    console.log("Deleting from list=>", photo);
    list.splice(index);
  }
}

function mergeEvents(current_list, in_list) {
  const photos_to_insert = []
  const photos_to_delete = []
  for (const photo of in_list.slice().reverse()) {
    if (photo.event != "deletion") {
      photos_to_insert.push(photo);
    } else {
      photos_to_delete.push(photo);
    }
  }
  for (const photo of photos_to_delete) {
    removePhotoFromList(current_list, photo);
    removePhotoFromList(photos_to_insert, photo);
  }
  for (const photo of photos_to_insert) {
    current_list.push(photo);
  }
}

export const photos = {
  namespaced: true,
  state: {
    all: { photos_list: [] },
    my: { photos_list: [] },
    last_timestamp: 0,
    status: {}
  },
  actions: {
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
    getRequest(state) {
      state.all.loading = true;
    },
    getSinceRequest(state) {
      state.all.loading = true;
    },
    deleteRequest(state) {
      state.all.loading = true;
    },
    getSinceSuccess(state, photos) {
      mergeEvents(state.all.photos_list, photos);
      if (state.all.photos_list.length > 0) {
        state.last_timestamp = state.all.photos_list[0].timestamp;
      }
      Vue.delete(state.all, "loading");
    },
    getOwnRequest(state) {
      state.all.loading = true;
    },
    getSuccess(state, photo) {
      if (photo.event != "deletion") {
        if (state.all.photos_list.length > 0) {
          var found = false;
          for (var i = 0; i < state.all.photos_list.length; i++) {
            if (state.all.photos_list[i].photo_id == photo.photo_id) {
              state.all.photos_list[i] = photo;
              found = true;
            }
          }
          if (!found) {
            state.all.photos_list.unshift(photo);
          }
        } else {
          var photos_list = [photo];
          state.all = { photos_list };
        }
      } else {
        removePhotoFromList(state.all.photos_list, photo);
      }
        
      Vue.delete(state.all, "loading");
    },
    getOwnSuccess(state, photos) {
      mergeEvents(state.my.photos_list, photos);
      Vue.delete(state.all, "loading");
    },
    deleteSuccess(state, photo) {
      if (typeof state.all.photos_list !== "undefined") {
        for (var i = 0; i < state.all.photos_list.length; i++) {
          if (state.all.photos_list[i].photo_id == photo.photo_id) {
            state.all.photos_list.splice(i, 1);
          }
        }
      }
      if (typeof state.my.photos_list !== "undefined") {
        for (i = 0; i < state.my.photos_list.length; i++) {
          if (state.my.photos_list[i].photo_id == photo.photo_id) {
            state.my.photos_list.splice(i, 1);
          }
        }
      }
      Vue.delete(state.all, "loading");
    },
    prepareEdit(state, id) {
      for (var i = 0; i < state.my.photos_list.length; i++) {
        if (state.my.photos_list[i].photo_id == id) {
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
        if (state.my.photos_list[i].photo_id == photo.photo_id) {
          Vue.set(state.my.photos_list, i, photo);
        }
      }
      if (
        typeof state.all.photos_list !== "undefined" &&
        state.all.photos_list.length > 0
      ) {
        for (i = 0; i < state.all.photos_list.length; i++) {
          if (state.all.photos_list[i].photo_id == photo.photo_id) {
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
        if (state.my.photos_list[i].photo_id == id) {
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
    getSinceFailure(state, error) {
      state.all = { error };
    },
    getOwnFailure(state, error) {
      state.all = { error };
    },
    editFailure(state, error) {
      state.all = { error };
    },
    getFailure(state, error) {
      if (error.error == "Item not found.") {
        for (var i = 0; i < state.all.photos_list.length; i++) {
          if (state.all.photos_list[i].photo_id == error.id) {
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
