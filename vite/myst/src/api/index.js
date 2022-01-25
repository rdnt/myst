import axios from "axios";

axios.defaults.baseURL = "/api";
// Add a request interceptor
axios.interceptors.request.use(
  function (config) {
    // Enable credentials passing (Authentication bearer token or cookie)
    config.withCredentials = true;
    // Tell the server we expect a JSON response
    config.headers = {
      "Content-type": "application/json"
    };
    return config;
  },
  function (error) {
    // Request failed, return a JSON error object
    return {
      status: "error",
      code: "REQUEST_INITIALIZATION_FAILED",
      message: error.message
    };
  }
);
// Add a response interceptor
// axios.interceptors.response.use(
//   function(response) {
//     // Verify response is a valid JSON encoded object
//     // if (typeof response.data !== "object") {
//     //   return {
//     //     status: "error",
//     //     code: "JSON_PARSE_FAILED",
//     //     message: "Response body is not a valid JSON object",
//     //     headers: response.headers
//     //   };
//     // }
//     // attach headers response data
//     // const headers = response.headers;
//     // response.headers = headers;
//     // return response;
//   },
//   function(error) {
//     // Catch failed requests, return JSON error string with REQUEST_FAILED code
//     if (error instanceof Error) {
//       return {
//         status: "error",
//         code: "REQUEST_FAILED",
//         message: error.message,
//         headers: {}
//       };
//     }
//   }
// );

export default axios;
