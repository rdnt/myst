import Vue from "vue";
import App from "@/App";
import WebSocket from "vue-native-websocket";
import store from "@/store";
import router from "@/router";

if (!WebAssembly.instantiateStreaming) {
  // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}
const go = new window.Go();
let mod, inst;
WebAssembly.instantiateStreaming(
  fetch("/assets/wasm/main0.wasm"),
  go.importObject
).then(result => {
  mod = result.module;
  inst = result.instance;
  go.run(inst);
  WebAssembly.instantiate(mod, go.importObject); // reset instance
});

Vue.prototype.$wasm = window.Go;

// import VueWasm from "./wasm";
// import Go from "@/../../public/assets/wasm/wasm_exec";
// console.log(Go);
// import go from "@/assets/wasm/main.wasm";

// import go from "@/../../public/assets/wasm/main.wasm";

// import wasmCode from "./lib.wasm";
//
// WebAssembly.compile(go).then(module => {
//   const instance = new WebAssembly.Instance(module);
//   console(instance.exports.add(1, 2)); // 3
// });

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

Vue.use(WebSocket, "ws://" + window.location.host, {
  reconnection: false,
  store: store,
  format: "json"
});

const init = async () => {
  // await VueWasm(Vue, { modules: { go } });
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
