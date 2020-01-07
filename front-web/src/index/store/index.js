import Vue from "vue"
import Vuex from "vuex"
import oidc from "./modules/oidc"
import zerodb from "./modules/zerodb"

Vue.use(Vuex)

const store = new Vuex.Store({
	state: {},
	modules: {
		zerodb,
		oidc
	},
	mutations: {},
	actions: {},
	getters: {}
})

export default store
