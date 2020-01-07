import {
	API_BASE_URL
} from "@env";
import axios from "axios";
import handleError from "@utils/catchError";
// import { getUserInfo } from './getUserToken'
import {
	loading
} from "./loading"
// import cookie from '@2dfire/share/Cookie'
const Loading = new loading()
// let token_type = getUserInfo().token_type || "Bearer";
// let access_token = getUserInfo().access_token;

// 超时时间
axios.defaults.timeout = 30000;
axios.defaults.withCredentials = true;
// http请求拦截器
axios.interceptors.request.use(
	config => {
		// config.headers.common["Authorization"] = token_type + " " + access_token;
		config.headers.common["Content-Type"] = "application/json;charset=utf-8";
		config.url = API_BASE_URL + config.url;
		config.params = Object.assign({}, config.params);
		//loading 默认为false 此处为了方便 全局设置为true
		config.loading === false ? config.loading = false : config.loading = true
		if (config.loading) {
			Loading.loadingStart(!config.loading || false)
		}
		return config;
	},
	error => {
		return Promise.reject(error);
	}
);
// http响应拦截器
axios.interceptors.response.use(
	res => {
		let data = res.data;
		if (data.code === 1) {
			Loading.loadingEnd(loading || false)
			return data.data;
		} else {
			Loading.loadingEnd(loading || false)
			handleError(res);
			return Promise.reject(res);
		}
	},
	error => {
		console.log(error.response, 'error11111111')
		Loading.loadingEnd(loading || false)
		handleError(error.response);
		return Promise.reject(error);
	},
);

export default axios;
