import Vue from 'vue';
import Vuex from 'vuex';
Vue.use(Vuex);

const state = {
    Wallet: "unknown",
    Accounts: [],
    BtpdList: {},
};
const store = new Vuex.Store({
    state
});

export default store;