const path = require("path");

module.exports = {
  // entry: ["@babel/polyfill", "src"],
  outputDir: path.resolve(__dirname, "../static/"),
  chainWebpack: config => {
    // config.optimization.delete('splitChunks')
    config.resolve.symlinks(false);
    config.module
      .rule("eslint")
      .use("eslint-loader")
      .options({
        fix: true
      });

    config.experiments = { asyncWebAssembly: true, importAsync: true };

    // config.module
    //  .rule('wasm')
    //    .test(/.wasm$/)
    //    .use('wasm-loader')
    //    .loader('wasm-loader')

    //  test: /\.wasm$/,
    // loaders: ['wasm-loader']
  },
  devServer: {
    public: "0.0.0.0:8081",
    host: "",
    port: 8081,
    contentBase: path.resolve(__dirname, "../assets/"),
    contentBasePublicPath: "/assets/",
    // publicPath: "./assets/",
    proxy: {
      "/api": {
        target: "http://0.0.0.0:8080",
        ws: true,
        changeOrigin: true,
        hostRewrite: true
      }
      // '/assets': {
      //   target: 'http://0.0.0.0:8081/',
      //   // changeOrigin: true,
      //   // hostRewrite: true,
      // },
    }
  },
  pluginOptions: {
    "style-resources-loader": {
      preProcessor: "scss",
      patterns: [path.resolve(__dirname, "./src/styles/App.scss")]
    }
  },
  css: {
    sourceMap: true
  }
};
