// import sites from '../../api/sites'
// import Vue from 'vue'
// initial state
const state = {
  all: []
};

// getters
const getters = {};

// actions
const actions = {
  // getSites({ commit }) {
  //     // sites.getSites(sites => {
  //     //Vue.prototype.$socket.send('test123')
  //     console.log(Vue.prototype.$socket.get('sites'))
  //     //commit('setSites', [])
  //     // })
  // }
};

// mutations
const mutations = {
  setSites(state, sites) {
    state.all = sites;
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};
