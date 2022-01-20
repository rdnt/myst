/* eslint-disable */
const path = require('path')

module.exports = {
  outputDir: path.resolve(__dirname, '../static/'),
  devServer: {
    public: '0.0.0.0:8082',
    host: '',
    port: 8082,
    proxy: {
      '/api': {
        target: 'http://0.0.0.0:8081',
        ws: true,
        changeOrigin: true,
        hostRewrite: true
      }
    }
  }
}
