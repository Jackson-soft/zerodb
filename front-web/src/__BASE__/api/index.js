// import axios from "@2dfire/share/Http";
import axios from "@2dfire/share/Http"
import qs from "@2dfire/share/QueryString"
import { BASE_URL } from "@env"

console.log("===BASE_URL", axios)
axios.setConfig({
	loading: false,
	isToken: false
})

// axios.defaults.headers['Content-Type'] = 'application/x-www-form-urlencoded;charset=utf-8';

const API = {
	SnapshotList(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/snapshot_list",
			params
		})
	},
	ClusterList(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/cluster_list",
			params
		})
	},
	Basic(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/basic",
			params
		})
	},
	creatSnapshot(params) {
		return axios({
			headers: {
				"Content-Type": "application/x-www-form-urlencoded;charset=utf-8"
			},
			method: "POST",
			url: BASE_URL + "/api/snapshot_config",
			data: qs.stringify(params)
		})
	},
	pushConfig(params) {
		return axios({
			headers: {
				"Content-Type": "application/x-www-form-urlencoded;charset=utf-8"
			},
			method: "POST",
			url: BASE_URL + "/api/push_config",
			data: qs.stringify(params)
		})
	},
	initConfig(params) {
		return axios({
			method: "POST",
			url: BASE_URL + "/api/config/init",
			data: params
		})
	},
	exportConfig(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/export",
			params
		})
	},
	stopService(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/stop_service",
			params
		})
	},
	switchConfig(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/switch",
			params
		})
	},
	hostgroups(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/hostgroups",
			params
		})
	},
	delHostgroups(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/delete_hostgroups",
			data: params
		})
	},
	addHostgroups(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/add_hostgroups",
			data: params
		})
	},
	hostgroupClusters(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/host_group_clusters",
			params
		})
	},
	addHostgroupcluster(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/add_hostgroupcluster",
			data: params
		})
	},
	delHostgroupcluster(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/delete_hostgroupcluster",
			data: params
		})
	},
	full(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/schema_configs",
			params
		})
	},
	addTable(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/add_table",
			data: params
		})
	},
	delTable(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/delete_table",
			data: params
		})
	},
	addSchema(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/add_schema",
			data: params
		})
	},
	delSchema(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/delete_schema",
			data: params
		})
	},
	updateShardrw(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/update_shardrw",
			data: params
		})
	},
	proxyStatus(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/proxy_cluster/status",
			params
		})
	},
	switchProxy(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/proxy_cluster/switch",
			data: params
		})
	},
	unregister(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/proxy/unregister",
			data: params
		})
	},
	updateStopservice(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/update_stopservice",
			data: params
		})
	},
	updateHostGroupCluster(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/update_hostgroupcluster",
			data: params
		})
	},
	updateBasic(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/update_basic",
			data: params
		})
	},
	updateSwitch(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/update_switch",
			data: params
		})
	},
	rollbackConfig(params) {
		return axios({
			headers: {
				"Content-Type": "application/json;charset=UTF-8"
			},
			method: "POST",
			url: BASE_URL + "/api/rollback_config",
			data: params
		})
	},
	getConfigTable(params) {
		return axios({
			method: "GET",
			url: BASE_URL + "/api/config/tables",
			params
		})
	}
}
export default API
