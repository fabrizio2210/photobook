module.exports = {
  devServer: {
    proxy: {
      "/api": { target: "http://localhost:5000" },
      "/static": {
        target: "http://localhost:80",
        pathRewrite: { "^/static": "" }
      }
    },
    watchOptions: {
      poll: 1000
    }
  }
};
