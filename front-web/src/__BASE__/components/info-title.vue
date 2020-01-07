<!--
    基础信息title zipai at 2018/6/4
    @params :
        title :标题,
-->
<template>
    <div class="info_title">
        <h4>{{title}}</h4>
    </div>
</template>
<script>
export default {
    props: ['title'],
};
</script>
<style lang="less" scoped>
.info_title {
    font-weight: bold;
    height: 44px;
    text-align: center;
    line-height: 44px;
    font-size: 14px;
    color: #333333;
}
</style>


