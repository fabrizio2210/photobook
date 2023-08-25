<template>
  <div class="credentialmanager">
    <div v-if="auth.user">
      <p>Logged in as: {{ auth.user.uid }}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "CredentialManager",
  data: function() {
    return {};
  },

  computed: {
    auth() {
      return this.$store.state.authentication;
    }
  },

  methods: {
    async getUid() {
      this.$store.dispatch("authentication/get_uid");
    },
    handlePhotoEvents(msg) {
      const evento = JSON.parse(msg);
      console.log("SSE:", evento);
      this.$store.dispatch("photos/mergeEvent", {evento});
    },
    populatePhotos() {
      this.$store.dispatch("photos/getAll");
    },
    setUid(id) {
      this.$store.dispatch("authentication/set_uid", {id});
    }
  },
  mounted() {
    if (typeof this.$route.query.id != "undefined") {
      this.setUid(this.$route.query.id);
    }
    if (!this.auth.user || !this.auth.user.uid) {
      this.getUid();
    }
    this.populatePhotos();
    this.$sse
      .create("/api/notifications")
      .on("photo", this.handlePhotoEvents)
      .on("error", err =>
        console.error("Failed to parse or lost connection:", err)
      )
      .connect()
      .catch(err => console.error("Failed make initial connection:", err));
  },
  sse: {
    cleanup: true
  }
};
</script>
