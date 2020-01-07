<script>
/**
     -- zipai@2dfire
     静态的显示按钮
     @props:
		 type:类型，显示不同的颜色

		 pending：办理中
		 resolved：已完结
*/

export default {
  name: 'ButtonStatic',
  props: ['type'],
  data() {
    return {}
  },
  computed: {
    getBg() {
      let style = ''
      switch (this.type) {
        case 'pending':
          style = 'red'
          break
        case 'resolved':
          style = 'green'
          break
        default:
          style = 'red'
          break
      }
      return style
    }
  }
}
</script>

<template>
    <div :class="['button',getBg]">
        <slot></slot>
    </div>
</template>

<style lang="less" scoped>
.button {
  vertical-align: middle;
  font-size: 10px;
  height: 20px;
  min-width: 50px;
  line-height: 20px;
  border-radius: 10px;
  background: #d83f3f;
  color: #ffffff;
  display: inline-block;
  padding: 0 10px;
}
.red{
  background: #d83f3f;
}
.green{
    background: #49C070;
}
</style>
