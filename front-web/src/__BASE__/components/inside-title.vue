<!--
    基础标题组件 zipai at 2018/08/08
    padding：15px;
-->
<template>
  <div class="inside-title">
    {{text}}
  </div>
</template>
<script>
export default {
  props: ['text']
}
</script>
<style lang="less" scoped>
  .inside-title{
    line-height: 25px;
    font-weight: bold;
    font-size: 18px;
    color: #3C4144;
  }
</style>
