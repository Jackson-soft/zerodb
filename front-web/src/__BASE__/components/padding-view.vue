<!--
    视图组件 zipai at 2018/08/07
    padding：15px;
-->
<template>
  <div class="padding15">
    <slot></slot>
  </div>
</template>
<script>
export default {}
</script>
<style lang="less" scoped>
.padding15 {
    background: #ffffff;
    padding: 15px;
    background: #ffffff;
    box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.13);
    margin-bottom: 20px;
}
</style>
