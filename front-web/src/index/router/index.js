import Vue from "vue"
import VueRouter from "vue-router"

Vue.use(VueRouter)
export default new VueRouter({
	routes: [
		{
			path: "*",
			redirect: "/flow_configure"
		},
		{
			path: "/oidc-callback",
			component: () => import("../views/login/oidc-callback"),
			meta: {
				isOidcCallback: true,
				isPublic: true
			},
			hidden: true
		},
		{
			path: "/oidc-callback-error",
			component: () => import("../views/login/oidc-callback-error"),
			hidden: true,
			meta: {
				isPublic: true
			}
		},
		{
			path: "/flow_configure",
			name: "配置",
			component: () => import("../views/configure")
		},
		{
			path: "/flow_colony",
			name: "集群",
			component: () => import("../views/colony")
		}
	],
	scrollBehavior() {
		// 路由变换后，滚动至顶
		return { x: 0, y: 0 }
	}
})
