<template>
  <div class="credentialmanager">
    <div v-if="auth.user">
      <p>Logged in as: {{ auth.user.uid }}</p>
    </div>
    <div class="closed-connection" v-if="closed_connection">
      <p>Connection lost with the server, refresh the page</p>
    </div>
    <div class="connection-error" v-else-if="connection_error">
      <p>Having connection issues. Reconnecting...</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "CredentialManager",
  data: function() {
    return {
      sse_client: {},
      closed_connection: false,
      connection_error: false
      };
  },

  computed: {
    auth() {
      return this.$store.state.authentication;
    }
  },

  methods: {
    handleReconnection(msg) {
      this.connection_error = false;
      this.closed_connetcion = false;
      this.$store.dispatch("photos/getAll");
      console.log("Opening:", msg);
    },
    handleFailedconnection(err) {
      this.closed_connetcion = true;
      console.error("Failed make initial connection:", err);
    },
    handlePhotoEvents(msg) {
      const evento = JSON.parse(msg);
      console.log("SSE:", evento);
      this.$store.dispatch("photos/mergeEvent", {evento});
    },
    handleUploadError(msg) {
      console.log("Error received from SSE:", msg);
      this.$store.dispatch("photos/setError", {msg});
    },
    handleConnectionError(err) {
      console.log("Failed to parse or lost connection:", err);
      console.log("readyState:", this.sse_client.source.readyState);
      this.connection_error = true;
      if (this.sse_client.source.readyState == 2) {
        this.closed_connection = true;
      }
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
      .on("error", this.handleConnectionError)
      .on("open", this.handleReconnection)
      .connect()
      .catch(this.handleFailedconnection)
      .then(client=>{ this.sse_client = client});
  },
  sse: {
    cleanup: true
  }
};
</script>
