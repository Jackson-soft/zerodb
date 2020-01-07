<!--
    视图组件  zipai at 2018/08/07
    小卡片
-->
<template>
    <div class="card-wrap">
        <slot></slot>
    </div>
</template>
<script>
export default {}
</script>
<style lang="less" scoped>
.card-wrap {
    display: inline-block;
    background: #ffffff;
    width: 200px;
    height: 100px;
    padding: 25px 15px;
    margin-bottom: 15px;
    background: #ffffff;
    border: 1px solid #e6e6e6;
    border-radius: 8px;
    margin-right: 35px;
}
</style>
