<template>
    <div id="index" class="index">
        <Nav :nav-data="navData" :nickName="nickName">
            <Icon v-if="isAdmin" type="gear-b" class="setting cu" @click.native="setting"></Icon>
        </Nav>
        <section class="main-wrap">
            <SideBar :list="list"></SideBar>
            <div class="right-wrap">
                <router-view></router-view>
            </div>
        </section>
    </div>
</template>

<script>
import Nav from "@components/nav/fui-nav"
import SideBar from "@components/side-bar/sidebar"
import list from "@components/side-bar/menu"

export default {
    name: "Index",
    data() {
        return {
            navData: {
                title: "ZeroDB",
                logo: "https://assets.2dfire.com/frontend/11dfef885ad8f93b77397febcdd8e127.png"
            },
            nickName: "--",
            isAdmin: false,
            list
        }
    },
    components: {
        Nav,
        SideBar
    },
    methods: {
        init() {
            this.Account()
        },
        setting() {
            window.location.href = "./manage.html"
            console.log("点击了设置")
        },
        Account() {
            // let self = this
            // API.Account()
            //   .then(data => {
            // 	self.nickName = data.full_name,
            // 	self.isAdmin = data.privileges.indexOf("ROLE_WORKFLOW_ADMIN")>-1
            //   })
            //   .catch(e => {})
        }
    },
    created() {
        this.init()
    }
}
</script>

<style lang="less" scoped>
#index {
    min-height: 100%;
    background: #f6f6f6;
}
.main-wrap {
    width: 100%;
    overflow: hidden;
    overflow-y: auto;
    padding-top: 60px;
    // height: calc(100% - 60px)
}
.right-wrap {
    margin-left: 200px;
    // height: 100%;
    // padding-top:60px;
}
.setting {
    margin-left: 20px;
    font-size: 24px;
    margin-top: 14px;
    color: #ffffff;
}
</style>
