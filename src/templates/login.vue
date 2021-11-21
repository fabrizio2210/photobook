<div id="app">

  <div>
    <p>Username:</p>
    <input v-model="username" placeholder="fabrizio2">
    <p>Password:</p>
    <input v-model="password" placeholder="pwd2">
    <button v-on:click="login()">Login</button>
    <button v-on:click="logout()">Logout</button>
    <p class="errorText" v-if="login_status == 2">{{ login_error }}</p>
  </div>

  <div v-if="login_status == 1">
    <a href="/web/projects">Your projects</a>
  </div>
  <div>
    <p>User ID: {{jwtData.identity}}</p>
    <p>Access: {{access_token}}</p>
    <p>Refresh: {{refresh_token}}</p>
  </div>

</div> 

<link rel="stylesheet" href="/css/stylesheet.css">
<script src="https://cdn.jsdelivr.net/npm/vue@2.5.13/dist/vue.js"></script>
<script src="/js/main.js"></script>
