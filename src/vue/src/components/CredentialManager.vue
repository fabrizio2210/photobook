<template>
  <div class="credentialmanager">
    <p class="errorText" v-if="login_status == 2">{{ login_error }}</p>
    <a v-if="login_status == 1" href="/web/projects">Your projects</a>
    <div v-if="auth.user" >
      <p>Logged in as: {{auth.username}}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "CredentialManager",
  data: function() {
    return {
      jwtData: {},
      access_token: "",
      refresh_token: "",
      login_status: 0,
      login_error: "",
    }
  }, 

  computed: {
    auth () {
      return this.$store.state.authentication;
    },
  },

  methods: {
    async refreshJWT() {
      if (this.auth.user) {
        if (this.auth.user.refresh_token) {
          const refresh_token = this.auth.user.refresh_token
          this.$store.dispatch('authentication/token_refresh',  refresh_token );
        }
      }
    },
  },
  mounted() {
    setInterval(this.refreshJWT, 20000);
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
