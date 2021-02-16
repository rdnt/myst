import Vue from "vue";
import Vuex from "vuex";
import * as util from "@/util";
import sites from "./modules/sites";
// import { v4 as uuidv4 } from "uuid";

const state = {
  loggedIn: false,
  socket: {
    isConnected: false,
    message: "",
    reconnectError: false
  },
  supports: {
    webAssembly: util.supportsWebAssembly
  }
};

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== "production";

export default new Vuex.Store({
  namespaced: true,
  modules: {
    sites
  },
  state,
  strict: debug,
  actions: {
    // get: function(_, type) {
    //   var id = Math.random()
    //     .toString()
    //     .slice(2, 10);
    //   // id = uuidv4();
    //   // var msg = {
    //   //   id: id,
    //   //   action: "get",
    //   //   type: type
    //   // };
    //   // Vue.prototype.$socket.send(JSON.stringify(msg));
    //   // Vue.prototype.$socket.onmessage = message => {
    //   //     if (message) {
    //   //         try {
    //   //             var response = JSON.parse(message.data)
    //   //         } catch (e) {
    //   //             console.log(e)
    //   //         }
    //   //     }
    //   //     // if (response['id'] == id) {
    //   //     //     return response
    //   //     // } else {
    //   //     //     console.log('response ignored')
    //   //     // }
    //   // }
    //   Vue.prototype.$socket.onmessage = message => {
    //     if (message) {
    //       try {
    //         var response = JSON.parse(message.data);
    //       } catch (error) {
    //         // console.log(error)
    //       }
    //       if (response["id"] === id) {
    //         this.commit("sites/setSites", response.data);
    //       } else {
    //         // console.log('response ignored')
    //       }
    //     }
    //   };
    // }
  },
  mutations: {
    SOCKET_ONOPEN(state, event) {
      Vue.prototype.$socket = event.currentTarget;
      state.socket.isConnected = true;
      // this.$store.dispatch('sites/getSites')
    },
    SOCKET_ONCLOSE(state) {
      state.socket.isConnected = false;
    },
    // SOCKET_ONERROR(state, event) {
    //     console.error(state, event)
    // },
    // default handler called for all methods
    SOCKET_ONMESSAGE(state, message) {
      state.socket.message = message;
    },
    // mutations for reconnect methods
    // SOCKET_RECONNECT(state, count) {
    //     console.info(state, count)
    // },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },

    login(state) {
      state.loggedIn = true;
    }
  }
});
