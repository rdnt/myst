import Vue from "vue";
import Vuex from "vuex";
import keystore from "@/store/modules/keystore";

Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    keystore
  }
});
