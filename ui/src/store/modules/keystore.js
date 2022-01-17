import api from "@/api";

export default {
  namespaced: true,
  state: {
    keystore: null
  },
  actions: {
    authenticate({ commit }, { keystoreId, passphrase }) {
      return api
        .post(`/keystore/${keystoreId}`, { passphrase })
        .then(response => {
          console.log(response);
          if (response.status === 200) {
            commit("setKeystore", response.data);
          }
        });
    },
    createEntry({ commit }, { keystoreId, entry }) {
      return api
        .post(`/keystore/${keystoreId}/entries`, entry)
        .then(response => {
          console.log(response);
          if (response.status === 200) {
            commit("setKeystore", response.data);
          }
        });
    }
  },
  mutations: {
    setKeystore(state, keystore) {
      keystore.entries = keystore.entries.slice(0, 40);
      state.keystore = keystore;
    }
  }
};
