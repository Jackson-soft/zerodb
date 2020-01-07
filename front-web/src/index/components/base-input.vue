<!--
    基础输入框 zipai at 2018/5/15
    @params :
        title :标题,
        type: 1 显示必选小红点 ，默认不显示必选小红点,
        placeholder:placeholder,
-->
<template>
    <div class="select-wrap clearfix">
        <div :class="['left',{'mdwidth':labelWidth==='md'}]">
            <OpenLabel :text="title" :type="type"></OpenLabel>
        </div>
        <div class="right">
            <slot name='input'></slot>
        </div>
    </div>
</template>
<script>
import OpenLabel from "@components/open-label";
export default {
    props: ["title", "type", "placeholder",'labelWidth'],
    data() {
        return {
            value: ""
        };
    },
    components: {
       OpenLabel,
    },
    methods: {
    },
    created() {}
};
</script>
<style lang="less" scoped>
.select-wrap {
    // overflow: hidden;
    margin-bottom: 20px;
}
.left {
    float: left;
    height: 30px;
    width: 120px;
    text-align: right;
    margin-right: 10px;
}
.mdwidth{
    width:94px;
}
.right {
    float: left;
    width: 310px;
}
</style>


