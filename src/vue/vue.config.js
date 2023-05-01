module.exports = {
  configureWebpack: {
    devServer: {
      disableHostCheck: true,
      watchOptions: {
        poll: 1000
      },
      port: 80,
      sockPort: 80,
    },
  },
};
