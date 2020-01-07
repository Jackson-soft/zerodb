<!--
    侧边栏组件 zipai at 2018/08/07

    传入list，自动生成侧边栏并分析跳转。

    @params :
        list: [
        {
          name: '/mnangeMent',
          title: '招商管理',
          icon: 'ios-paper',
          list: [
            {
              name: '/mnangeMent_brand',
              icon: 'ios-paper',
              title: '品牌管理',
              list: [
                {
                  name: '/mnangeMent_brand_infoMaintain',
                  title: '品牌信息维护'
                },
                {
                  name: '/mnangeMent_brand_progressControl',
                  title: '品牌进度管理'
                },
                {
                  name: '/mnangeMent_brand_infoApproval',
                  title: '品牌信息审批'
                }
              ]
            },
            {
              name: '/mnangeMent_assignment',
              icon: 'ios-paper',
              title: '招商任务管理'
            }
          ]
        }
      ]
      支持三种方式侧边栏项：1. 有二级菜单 2.仅有一级菜单 3.菜单直接跳转

-->
<template>
	<!-- eslint-disable -->
	<div class="sb-wrap">
		<Menu
			ref="leftMenu"
			theme="dark"
			@on-select="onselect"
			:active-name="activeName"
			:open-names="openNames"
			width="200px"
		>
			<div v-for="(item, index) in list" :key="index">
				<div v-if="item.list">
					<Submenu :name="item.name">
						<template slot="title">
							<Icon :type="item.icon"></Icon>{{ item.title }}
						</template>
						<div v-for="(secondItem, index) in item.list" :key="index">
							<Submenu v-if="secondItem.list" :name="secondItem.name">
								<template slot="title">
									<Icon :type="secondItem.icon"></Icon>{{ secondItem.title }}
								</template>
								<MenuItem v-for="(thirdItem, index) in secondItem.list" :key="index" :name="thirdItem.name">{{
									thirdItem.title
								}}</MenuItem>
							</Submenu>
							<MenuItem v-else :name="secondItem.name" :key="index"> {{ secondItem.title }}</MenuItem>
						</div>
					</Submenu>
				</div>
				<div v-else>
					<MenuItem :name="item.name">{{ item.title }}</MenuItem>
				</div>
			</div>
		</Menu>
	</div>
</template>
<script>
const nameArr = []
function getActiveName(data) {
	data.forEach(v => {
		nameArr.push(v.name)
		if (v.list) {
			getActiveName(v.list)
		}
	})
}
export default {
	name: 'SideBar',
	props: ['list'],
	data() {
		return {
			activeName: '',
			openNames: []
		}
	},
	methods: {
		onselect(path) {
			this.$router.push({
				path
			})
		},
		getRoutePath() {
			const currentPath = this.$route.path
			const routerArr = currentPath.split('_')
			const firstOrder = routerArr[0]
			const secondOrder = routerArr.slice(0, 2).join('_')
			const thrOrder = routerArr.slice(0, 3).join('_')
			if (nameArr.includes(thrOrder)) {
				this.openNames = [].concat([firstOrder, secondOrder])
				this.activeName = thrOrder
			} else {
				this.openNames = [].concat([firstOrder])
				this.activeName = secondOrder
			}
			this.$nextTick(() => {
				this.$refs.leftMenu.updateOpened()
				this.$refs.leftMenu.updateActiveName()
			})
		}
	},
	watch: {
		$route() {
			this.getRoutePath()
		}
	},
	created() {
		getActiveName(this.list)
	},
	mounted() {}
}
</script>
<style lang="less" scoped>
.sb-wrap {
	position: fixed;
	width: 200px;
	height: 100%;
	top: 60px;
	left: 0;
	overflow: auto;
	overflow-x: hidden;
	overflow-y: hidden;
	background-color: #333333;
}
</style>
