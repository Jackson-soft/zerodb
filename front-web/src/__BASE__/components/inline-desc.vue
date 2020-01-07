<script>
    /**
     -- zipai@2dfire

     @props:
     title：左侧label，
     desc：右侧文字，也可以插入slot任意行内元素，元素垂直对齐：top
     行内 "标签: 内容"展示，其中左侧标签宽度固定，右侧占剩余宽度
     */

    export default {
        props: ["title"],
        components: {},
        methods: {},
        created() {

        }
    }
</script>

<template>
    <div class="content-wrap">
        <span class="left">{{title}}</span>
        <span class="right">
            <slot></slot>
        </span>
    </div>
</template>

<style lang="less" scoped>
    .content-wrap {
        margin-top: 14px;
        font-size: 12px;
        width: 100%;
    }

    .left {
        display: inline-block;
        width: 100px;
        text-align: right;
        color: #999999;
    }

    .right {
        display: inline-block;
        padding-left: 30px;
        max-width: calc(100% - 130px);
        color: #333333;

        vertical-align: top;
    }
</style>

