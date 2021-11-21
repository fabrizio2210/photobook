

var app = new Vue({
  el: '#app',
  data: {
    username: "fabrizio2",
    password: "pwd2",
    jwtData: {},
    access_token: "",
    refresh_token: "",
    login_status: 0,
    login_error: "",
  }, 
  computed: {
  },
  methods: {
    async login() {
      // Error handling and such omitted here for simplicity.
      const res = await fetch(`/auth`,{
        method: 'POST',
        headers: new Headers({
          'Content-Type': 'application/json'
        }),
        body: JSON.stringify({"username": this.username, "password" : this.password})
      }).catch(error => {
        console.error('Proble during login:', error);
      });
      if (res.ok){
        var jwt = await res.json();
        this.jwtData.identity= JSON.parse(atob(jwt['access_token'].split('.')[1])).identity;
        this.access_token= await jwt['access_token'];
        this.refresh_token= jwt['refresh_token'];
        this.login_status = 1;
        this.login_error = 1;
      } else {
        this.login_status = 2;
        this.login_error = (await res.json())['message'];
      }
    },

    logout: function() {
      this.access_token = "";
      this.refresh_token = ""; 
    },

    async fetchJWT() {
      if (this.refresh_token) {
        // Error handling and such omitted here for simplicity.
        const res = await fetch(`/refresh`,{
          method: 'POST',
          headers: new Headers({
            'Authorization': `Bearer ${this.refresh_token}`
          }),
        });
        var jwt = await res.json();
        this.access_token= jwt['access_token'];
      }
    },
  },
  mounted() {
    setInterval(this.fetchJWT, 5000);
  }
})


