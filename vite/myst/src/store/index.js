import keystore from "./modules/keystore";
import { createStore } from 'vuex'

// Create a new store instance.
export default createStore({
    modules: {
        keystore
    },
    state () {
        return {
            ready: false,
            onboarding: false,
        }
    },
    mutations: {
        setReady (state, ready) {
            state.ready = ready;
        },
        setOnboarding (state, onboarding) {
            state.onboarding = onboarding;
        },
    }
})
