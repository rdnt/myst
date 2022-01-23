import api from "../../api";

export default {
    namespaced: true,
    state: {
        keystore: null,
        keystores: [],
        currentKeystore: null,
    },
    actions: {
        // authenticate({ commit }, { keystoreId, password }) {
        //   return api
        //     .post(`/keystore/${keystoreId}`, { password })
        //     .then(response => {
        //       console.log(response);
        //       if (response.status === 200) {
        //         commit("setKeystore", response.data);
        //       }
        //     });
        // },
        createEntry({commit}, {keystoreId, entry}) {
            return api
                .post(`/keystore/${keystoreId}/entries`, entry)
                .then(response => {
                    console.log(response);
                    if (response.status === 200) {
                        commit("setKeystore", response.data);
                    }
                });
        },
        getKeystores({commit}) {
            return api
                .get(`/keystores`)
                .then(response => {
                    console.log(response);
                    if (response.status === 200) {
                        commit("setKeystores", response.data);
                    }
                });
        }
    },
    mutations: {
        setKeystore(state, keystore) {
            keystore.entries = keystore.entries.slice(0, 40);
            state.keystore = keystore;
        },
        setKeystores(state, keystores) {
            state.keystores = keystores;
        }
    }
};
