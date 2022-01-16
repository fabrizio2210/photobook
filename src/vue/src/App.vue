<template>
  <div class="jumbotron" id="app">
    <div id="nav">
      <router-link to="/">Home</router-link> |
      <router-link to="/upload">Upload</router-link> |
      <router-link to="/edit">Edit</router-link>
    </div>
    <router-view />
    <CredentialManager />
    <div class="container">
      <div class="row">
        <div class="col-sm-6 offset-sm-3">
          <div v-if="alert.message" :class="`alert ${alert.type}`">
            {{ alert.message }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import CredentialManager from "@/components/CredentialManager.vue";

export default {
  components: {
    CredentialManager
  },
  computed: {
    alert() {
      return this.$store.state.alert;
    }
  },
  watch: {
    $route() {
      // clear alert on location change
      this.$store.dispatch("alert/clear");
    }
  }
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
