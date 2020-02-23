import Vue from "vue";
import App from "@/App";
import WebSocket from "vue-native-websocket";
import store from "@/store";
import router from "@/router";
import VueWasm from "./wasm";
import wasmExec from "@/assets/wasm/wasm_exec.js";
import go from "@/assets/wasm/main.wasm";

// import WasmExec from "@/../../public/assets/wasm/wasm_exec.js";
// import("@/../../public/assets/wasm/wasm_exec.js");
// import Go from "@/../../public/assets/wasm/main.wasm";
// import("@/../../public/assets/wasm/main.wasm");

// const extractModule = async module => {
//   const { instance } = await module();
//   return instance.exports;
// };
//
// const wasm = extractModule(Go);

// VueWasm(Vue, { modules: { go } });

Vue.config.productionTip = false;

// Vue.prototype.$wasm = wasm;

Vue.use(WebSocket, "ws://" + window.location.host, {
  reconnection: false,
  store: store,
  format: "json"
});

const init = async () => {
  await VueWasm(Vue, { modules: { go } });
  /* eslint-disable no-new */
  new Vue({
    store,
    router,
    render: h => h(App)
  }).$mount("#app");
};

init();

// new Vue({
//   store,
//   router,
//   render: h => h(App)
// }).$mount("#app");
