module.exports = {
  devServer: {
    proxy: {
      "^/api": {
        target: "http://localhost:9900",
        changeOrigin: true,
        logLevel: "debug",
      },
    },
  },
};
