import api from "../../api";

export default {
    namespaced: true,
    state: {
        keystore: null,
        keystoreIds: [],
        keystores: [],
        currentKeystore: null,
        onboarding: false,
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
        createEntry ({ commit }, { keystoreId, entry }) {
            return api
                .post(`/keystore/${keystoreId}/entries`, entry)
                .then(response => {
                    console.log(response);
                    if (response.status === 200) {
                        commit("setKeystore", response.data);
                    }
                });
        },
        createKeystore ({ commit }, { name, password }) {
            return api
                .post(`/keystores`, { name, password })
                .then(response => {
                    console.log(response);
                    if (response.status === 200) {
                        commit("setKeystore", response.data);
                    }
                });
        },
        getKeystoreIds ({ commit }) {
            return api
                .get(`/keystores`)
                .then(response => {
                    console.log(response);
                    if (response.status === 200) {
                        commit("setKeystoreIds", response.data);
                    }
                });
        }
    },
    mutations: {
        setKeystore (state, keystore) {
            keystore.entries = keystore.entries.slice(0, 40);
            state.keystore = keystore;
        },
        setKeystoreIds (state, ids) {
            state.keystoreIds = ids;
        },
        setKeystores (state, keystores) {
            state.keystores = keystores;
        },
    }
};
