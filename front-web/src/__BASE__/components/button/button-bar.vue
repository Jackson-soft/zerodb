<script>
/**
     -- zeqi@2dfire

     @props:
     type：类型

     放置操作按钮组 右置，下边距15px
     */

export default {
  props: [''],
  data () {
    return {
    }
  },
  methods: {

  },
  created () {

  }
}
</script>

<template>
<div class="right-buttons">
    <slot></slot>
</div>
</template>

<style lang="less" scoped>
    .right-buttons {
        width: 100%;
        text-align: right;
        margin-bottom: 15px;
    }

</style>
