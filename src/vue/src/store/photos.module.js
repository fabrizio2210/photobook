import { photoService } from "../services";
import Vue from "vue";

function removePhotoFromList(list, photo) {
  for (var index in list) {
    if (list[index].photo_id == photo.photo_id) {
      if (photo.timestamp > list[index].timestamp) {
        console.log("Deleting from list=>", photo);
        list.splice(index, 1);
        index--;
      }
    }
  }
}

function mergeEvents(current_list, in_list) {
  const photos_to_insert = [];
  const photos_to_delete = [];
  for (const photo of in_list.slice()) {
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
    var found = false;
    for (var i = 0; i < current_list.length; i++) {
      if (current_list[i].photo_id == photo.photo_id) {
        console.log("Found to substitute=>", current_list[i]);
        if (photo.timestamp > current_list[i].timestamp) {
          Vue.set(current_list, i, photo);
        }
        found = true;
        break;
      }
    }
    if (!found) {
      current_list.unshift(photo);
    }
  }
}

export const photos = {
  namespaced: true,
  state: {
    all: { photos_list: [] },
    my: { photos_list: [] },
    ticket_id: "",
    status: {}
  },
  actions: {
    prepareEdit({ commit }, { id }) {
      commit("prepareEdit", id);
    },
    edit({ commit }, { uid, photo }) {
      photoService.put(uid, photo).then(
        photo => commit("editSuccess", photo["event"]),
        error => commit("editFailure", error)
      );
    },
    unedit({ commit }, { id }) {
      commit("unedit", id);
    },
    getAll({ commit }) {
      commit("getAllRequest");
      photoService.getAll().then(
        photos => commit("getAllSuccess", photos["events"]),
        error => commit("getAllFailure", error)
      );
    },
    putMetadata({ commit, dispatch }, metadata) {
      photoService
        .putMetadata(metadata)
        .then(dispatch("getTicket"), error =>
          commit("putMetadataFailure", error)
        );
    },
    get({ commit }, { id }) {
      commit("getRequest");
      photoService.get(id).then(
        photo => commit("getSuccess", photo["event"]),
        error => commit("getFailure", { error, id })
      );
    },
    getTicket({ commit }) {
      photoService.getTicket().then(
        ticket => commit("getTicketSuccess", ticket["ticket_id"]),
        error => commit("getTicketFailure", error)
      );
    },
    getOwn({ commit }, { uid }) {
      commit("getOwnRequest");

      photoService.getOwn(uid).then(
        photos => commit("getOwnSuccess", photos["events"]),
        error => commit("getOwnFailure", error)
      );
    },
    del({ commit }, { uid, id }) {
      commit("deleteRequest");
      photoService.del(uid, id).then(
        photo => commit("deleteSuccess", photo["event"]),
        error => commit("deleteFailure", error)
      );
    },
    mergeEvent({ commit }, { evento }) {
      commit("mergePhotoEvents", [evento]);
    },
    resetError({ commit }) {
      commit("resetError");
    },
    setError({ commit }, { msg }) {
      commit("setError", msg);
    }
  },
  mutations: {
    getRequest(state) {
      state.all.loading = true;
      state.status = {};
    },
    getAllRequest(state) {
      state.all.loading = true;
      state.status = {};
    },
    deleteRequest(state) {
      state.all.loading = true;
      state.status = {};
    },
    getAllSuccess(state, photos) {
      state.status = {};
      mergeEvents(state.all.photos_list, photos);
      Vue.delete(state.all, "loading");
    },
    getOwnRequest(state) {
      state.all.loading = true;
      state.status = {};
    },
    getSuccess(state, photo) {
      state.status = {};
      mergeEvents(state.all.photos_list, [photo]);
      Vue.delete(state.all, "loading");
    },
    getTicketSuccess(state, ticket_id) {
      state.ticket_id = ticket_id;
    },
    mergePhotoEvents(state, events) {
      mergeEvents(state.all.photos_list, events);
      Vue.delete(state.all, "loading");
    },
    getOwnSuccess(state, photos) {
      state.status = {};
      mergeEvents(state.my.photos_list, photos);
      Vue.delete(state.all, "loading");
    },
    deleteSuccess(state, photo) {
      state.status = {};
      mergeEvents(state.all.photos_list, [photo]);
      mergeEvents(state.my.photos_list, [photo]);
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
      state.status = {};
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
    getAllFailure(state, error) {
      state.status = { error };
      state.all = { error };
    },
    getTicketFailure(state, error) {
      state.status = { error };
    },
    putMetadataFailure(state, error) {
      state.status = { error };
    },
    getOwnFailure(state, error) {
      state.status = { error };
      state.all = { error };
    },
    editFailure(state, error) {
      state.status = { error };
      state.all = { error };
    },
    resetError(state) {
      state.status = {};
    },
    setError(state, error) {
      state.status = { error };
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
        state.status = { error };
      }
    },
    deleteFailure(state, error) {
      state.all = { error };
      state.status = { error };
    }
  }
};
