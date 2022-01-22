import { createStore } from 'vuex'
import keystore from "./modules/keystore";

export default createStore({
  modules: {
    keystore
  }
})
