const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function (app) {
  app.use(
    createProxyMiddleware(["/github/**", "/api/**"], {
      target: "http://localhost:8000",
    })
  );
};
