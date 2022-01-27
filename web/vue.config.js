module.exports = {
  devServer: {
    proxy: {
      "^/api": {
        target: "http://localhost:9292",
        changeOrigin: true,
        logLevel: "debug",
      },
    },
  },
};
