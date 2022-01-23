import {createStore} from 'vuex'
import keystore from "./modules/keystore";

export default createStore({
    modules: {
        keystore
    },
    state: {
        ready: false
    },
    mutations: {
        setReady(state, ready) {
            state.ready = ready;
        },
    }
})
