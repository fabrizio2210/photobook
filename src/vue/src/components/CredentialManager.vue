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
    handlePhotoEvents(msg) {
      const evento = JSON.parse(msg);
      console.log("SSE:", evento);
      this.$store.dispatch("photos/mergeEvent", {evento});
    },
    handleUploadError(msg) {
      console.log("Error received from SSE:", msg);
      this.$store.dispatch("photos/setError", {msg});
    }
  },
  async mounted() {
    this.$store.dispatch("photos/getAll");
    if (typeof this.$route.query.id != "undefined") {
      const id = this.$route.query.id;
      this.$store.dispatch("authentication/set_uid", {id});
    }
    if (!this.auth.user || !this.auth.user.uid) {
       await this.$store.dispatch("authentication/get_uid");
    }
    this.$sse
      .create("/api/notifications/" + this.auth.user.uid)
      .on("photo", this.handlePhotoEvents)
      .on("error_upload", this.handleUploadError)
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
