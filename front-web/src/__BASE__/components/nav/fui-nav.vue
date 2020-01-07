<template>
    <!-- eslint-disable -->
    <header class="top-nav clearfix">
        <div class="logo fl cu" @click="toIndex">
            <img :src="getLogo" />
            <span class="title">{{ getTitle }}</span>
        </div>
        <div class="nav fr">
            <div class="fl">
                <Dropdown class="layout-dropdown" trigger="click">
                    <a class="my-account" href="javascript:void(0)">
                        <img class="photo" :src="oidcUser.picture" />
                        {{ oidcUser.nickname }}
                        <Icon class="icon-down-arrow" type="chevron-down"></Icon>
                    </a>
                    <Dropdown-menu slot="list">
                        <Dropdown-item name="exit" @click.native="logout">注销</Dropdown-item>
                    </Dropdown-menu>
                </Dropdown>
            </div>
            <div class="fl">
                <slot></slot>
            </div>
        </div>
    </header>
</template>
<script>
import API from "@api"
import { mapGetters } from "vuex"
// import { getNickName } from '@utils/getUserToken'
const demoData = {
    title: "管理后台",
    logo: "https://assets.2dfire.com/frontend/11dfef885ad8f93b77397febcdd8e127.png"
    // list: [
    //   {
    //     name: '供应链',
    //     val: 'supply',
    //     children: [
    //       {
    //         name: '库存',
    //         val: 'store'
    //       }
    //     ]
    //   },
    //   {
    //     name: '掌柜工具',
    //     val: 'boss'
    //   }
    // ]
}

export default {
    props: ["nav-data", "nickName"],
    data() {
        const data = this.navData || demoData
        return {
            state: data
        }
    },
    computed: {
        ...mapGetters("oidc", ["oidcUser"]),
        getLogo() {
            const defaultLogo = "https://assets.2dfire.com/frontend/11dfef885ad8f93b77397febcdd8e127.png"
            return this.state.logo || defaultLogo
        },
        getTitle() {
            const defaultTitle = "管理后台"
            return this.state.title || defaultTitle
        }
    },
    methods: {
        init() {},
        logout() {
            API.Logout().then(data => {
                console.log(data)
            })
        },
        // 回首页
        toIndex() {
            window.location.href = "./index.html"
        }
    },
    created() {
        // this.init()
    }
}
</script>

<style lang="less" scoped>
@theme-background-color: #33363a;
.top-nav {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 10;
    background-color: @theme-background-color;
    padding: 0 40px;
    height: 60px;
}
.logo {
    padding-top: 13px;
    img {
        width: 90px;
    }
    .title {
        display: inline-block;
        color: #fff;
        font-size: 16px;
        vertical-align: top;
        padding: 5px 0 0 10px;
        margin-left: 10px;
        position: relative;
        &:after {
            content: "";
            display: block;
            width: 1px;
            height: 20px;
            background-color: rgba(255, 255, 255, 0.6);
            position: absolute;
            left: 0;
            top: 6px;
        }
    }
}
.layout-dropdown {
    padding-top: 13px;
    margin-right: 15px;
}

.my-account {
    color: #fff;
    font-size: 14px;
}
.photo {
    display: inline-block;
    width: 30px;
    height: 30px;
    margin-right: 5px;
    vertical-align: middle;
    border-radius: 50%;
}
.icon-down-arrow {
    font-size: 12px;
    margin-left: 5px;
    vertical-align: middle;
}
.custom {
    float: right;
}
</style>
