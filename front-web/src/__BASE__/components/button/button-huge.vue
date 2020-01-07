<script>
/**
     -- zeqi@2dfire

     @props:
     width: 按钮宽度 默认内部内容撑开 左右15px内边距

     @events:
     on-tap: 点击的冒泡事件

     @slot:
     内容文字

     */

export default {
  name: 'ButtonHuge',
  props: {
    width: {
      type: [Number, String],
      default: 240
    },
    center: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {}
  },
  computed: {
    styles () {
      return {
        width: `${this.width}px`
      }
    },
    classes () {
      if (this.center) {
        return [
          'center'
        ]
      }
    }
  },
  methods: {
    emitTap () {
      this.$emit('on-tap')
    }
  }
}
</script>

<template>
    <div :class="classes">
    <button class="ivu-btn ivu-btn-primary button-huge" @click="emitTap" :style="styles">
        <slot></slot>
    </button>
    </div>
</template>

<style lang="less" scoped>
    .button-huge {
        padding: 0 15px;
        font-size: 14px;

        height: 44px;
        line-height: 44px;

        margin: 0 auto;
    }

    .center {
        text-align: center;
    }
</style>
