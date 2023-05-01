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
    }
  },
  mounted() {
    if (!this.auth.user || !this.auth.user.uid) {
      this.getUid();
    }
  }
};
</script>
