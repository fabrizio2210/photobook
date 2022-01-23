<template>
  <div>
    <em v-if="photos.loading">Loading photos...</em>
    <img v-show="photos.loading" src="../assets/loading.gif" />
    <span v-if="photos.error" class="text-danger"
      >ERROR: {{ photos.error }}</span
    >
    <div class="photo-grid" v-if="photos.photos_list">
      <div
        class="photo-row"
        v-for="photo in photos.photos_list"
        :key="photo.id"
      >
        <div class="btn-container">
          <button
            class="photo-col-btn btn-delete"
            v-show="!photo.edit"
            @click.prevent="deletePhoto(photo.id)"
          >
            Delete
          </button>
        </div>
        <div class="btn-container">
          <button
            class="photo-col-btn btn-edit"
            v-show="!photo.edit"
            @click.prevent="editPhoto(photo.id)"
          >
            Edit
          </button>
        </div>
        <div class="photo-col-img">
          <img
            class="photo-img-edit-element"
            loading="lazy"
            :src="photo.location"
          />
          <div v-show="!photo.edit" class="photo-description">
            {{ photo.description }}
          </div>
          <div v-show="!photo.edit && photo.author" class="photo-author">
            -- {{ photo.author }} --
          </div>
          <div v-show="photo.edit" class="text-input-group">
            <div class="form-group">
              <label
                >Description:
                <textarea
                  cols="40"
                  rows="5"
                  class="form-description"
                  v-model="photo.description"
                />
              </label>
            </div>
            <div class="form-group">
              <label
                >Name:
                <input type="text" class="form-author" v-model="photo.author" />
              </label>
              <button class="btn-delete" @click.prevent="uneditPhoto(photo.id)">
                Cancel
              </button>
              <button
                class="btn-ok"
                @click.prevent="confirmEditPhoto(photo.id)"
              >
                OK
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  },
  computed: {
    photos() {
      return this.$store.state.photos.my;
    },
    last_timestamp() {
      return this.$store.state.photos.last_timestamp;
    },
    uid() {
      if (this.$store.state.authentication.user) {
        if (this.$store.state.authentication.user.uid) {
          return this.$store.state.authentication.user.uid;
        }
      }
      return "12345";
    }
  },
  methods: {
    populatePhotos(uid) {
      if (typeof uid !== "undefined") {
        this.$store.dispatch("photos/getOwn", { uid });
      }
    },
    editPhoto(id) {
      this.$store.dispatch("photos/prepareEdit", { id });
    },
    uneditPhoto(id) {
      this.$store.dispatch("photos/unedit", { id });
    },
    confirmEditPhoto(id) {
      const photos = this.photos.photos_list;
      var photo;
      for (var i = 0; i < photos.length; i++) {
        if (photos[i].id == id) {
          photo = photos[i];
        }
      }
      const uid = this.uid;
      this.$store.dispatch("photos/edit", { uid, photo });
    },
    deletePhoto(id) {
      const uid = this.uid;
      this.$store.dispatch("photos/del", { uid, id });
    }
  },
  mounted() {
    this.populatePhotos(this.uid);
  },
  destroyed() {}
};
</script>
