<template>
  <div>
    <em v-if="photos.loading">Loading photos...</em>
    <img v-show="photos.loading" src="../assets/loading.gif" />
    <span v-if="photos.error" class="text-danger"
      >ERROR: {{ photos.error }}</span
    >
    <div class="photo-list" v-if="photos.photos_list">
      <div
        class="photo-element"
        v-for="photo in photos.photos_list"
        :key="photo.photo_id"
      >
        <img
          class="photo-img-element"
          loading="lazy"
          @error="imgError(photo.photo_id)"
          :src="photo.location"
        />
        <div class="photo-description">{{ photo.description }}</div>
        <div v-show="photo.author" class="photo-author">
          -- {{ photo.author }} --
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
    user() {
      return this.$store.state.authentication.user;
    },
    photos() {
      return this.$store.state.photos.all;
    },
    last_timestamp() {
      return this.$store.state.photos.last_timestamp;
    }
  },
  sse: {
    cleanup: true
  },
  methods: {
    populatePhotos(last_timestamp) {
      this.$store.dispatch("photos/getSince", { last_timestamp });
    },
    imgError(id) {
      this.$store.dispatch("photos/get", { id });
    },
    handlePhotoEvents(msg) {
      const evento = JSON.parse(msg);
      console.log("SSE:", evento);
      this.$store.dispatch("photos/mergeEvent", {evento});
    }
  },
  mounted() {
    this.populatePhotos(this.$store.state.photos.last_timestamp);
    this.vueInsomnia().on();
    this.$sse
      .create("/api/notifications")
      .on("photo", this.handlePhotoEvents)
      .on("error", err =>
        console.error("Failed to parse or lost connection:", err)
      )
      .connect()
      .catch(err => console.error("Failed make initial connection:", err));
  },
  destroyed() {
    this.vueInsomnia().off();
  }
};
</script>
