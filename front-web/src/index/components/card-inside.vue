<!--
    视图组件  zipai at 2018/08/07
    小卡片
-->
<template>
    <div class="card-inside clearfix">
        <img class="img fl" :src="imgUrl" alt="" />
        <div class="right fl">
            <span class="title">{{ text }}</span
            ><br />
            <span class="go cu" @click="emitTap">发起申请</span>
        </div>
    </div>
</template>
<script>
export default {
    props: ["img", "text"],
    data() {
        return {
            imgUrl: this.img ? this.img : "https://assets.2dfire.com/frontend/1f93549b3e857ad7a6923de6a1e7e63e.png"
        }
    },
    methods: {
        emitTap() {
            this.$emit("on-tap")
        }
    }
}
</script>
<style lang="less" scoped>
.card-inside {
    .img {
        width: 45px;
        height: 45px;
        margin-right: 20px;
    }
    .right {
        text-align: center;
    }
    .title {
        display: inline-block;
        width: 100px;
        text-overflow: ellipsis;
        white-space: nowrap;
        overflow: hidden;
        line-height: 28px;
        font-weight: bold;
        font-size: 14px;
        color: #333333;
    }
    .go {
        font-size: 12px;
        color: #5599ff;
    }
}
</style>
