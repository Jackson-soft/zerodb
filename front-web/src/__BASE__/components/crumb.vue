<!--
    面包屑组件 zipai at 2018/08/07

    自动分析路由，添加展示及跳转url

    在需要面包屑的页面直接引用组件并放在所需位置即可。
-->
<template>
	<!-- eslint-disable -->
	<div class="crumb-wrap">
		<span class="icon icon-crumb"></span>
		<Breadcrumb class="crumb-box" :separator="'/'">
			<Breadcrumb-item v-for="(item, index) in componentData.crumbs" :href="getURL(item, index)" :key="item.id">
				{{ getName(item) }}
			</Breadcrumb-item>
		</Breadcrumb>
		<div class="right-slot"><slot></slot></div>
	</div>
</template>
<script>
const componentData = {
	crumbs: []
}
const routes = []
export default {
	data() {
		return {
			componentData,
			routes
		}
	},
	methods: {
		getItem(name) {
			let item = {}
			this.routes.map(k => {
				if (k.path.split('/')[1] === name) {
					item = k
				}
				return item
			})

			return item
		},
		getURL(name, index) {
			let url = this.getItem(name).path
			if (this.$route.query && index > 0) {
				const query = '?' + this.$route.fullPath.split('?')[1]
				url += query
			}
			return url
		},
		getName(name) {
			return this.getItem(name).name
		},
		setRoutes() {
			const r = this.$router.options.routes
			console.log(r)
			this.routes = r
		},
		setCrumbs() {
			let { hash } = window.location
			const crumb = []

			if (hash.indexOf('?') > -1) {
				[hash] = hash.split('?')
			}
			/* eslint-disable */
			let routes = hash.split('/').splice(1)
			let RouteArr = routes[0].split('_')
			let RouteLength = RouteArr.length
			for (RouteLength; RouteLength > 1; RouteLength--) {
				crumb.unshift(RouteArr.slice(0, RouteLength).join('_'))
			}
			let routesArr = this.routes
			let routesPath = []
			routesArr.forEach(v => {
				routesPath.push(v.path)
			})
			if (!routesPath.includes(`/${crumb[0]}`)) {
				crumb.shift()
			}
			this.componentData.crumbs = crumb
		}
	},
	created() {
		this.setRoutes()
		this.setCrumbs()
	}
}
</script>

<style lang="less" scoped>
.crumb-wrap {
	width: 100%;
	height: 50px;
	line-height: 50px;
	// margin-bottom: 20px;
	padding-left: 15px;
	background: #ffffff;
	position: relative;
}

.crumb-box {
	color: #333333;
	font-size: 14px;
}

.icon {
	width: 3px;
	height: 18px;
	margin-right: 8px;
	float: left;
	margin-top: 14px;
	background: #d83f3f;
}
.right-slot {
	position: absolute;
	right: 10px;
	top: 0px;
}
</style>
