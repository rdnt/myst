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
    }
  },
  mutations: {
    setKeystore(state, keystore) {
      state.keystore = keystore;
    }
  }
};
