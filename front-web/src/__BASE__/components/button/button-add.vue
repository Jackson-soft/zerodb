<script>
/**
     -- zeqi@2dfire

     @props:

     @events:
     on-tap: 点击的冒泡事件

     @desc:
     新增腰圆按钮 左侧加号 右侧内容

     @slot:
     右侧插槽内容

     */

export default {
  name: 'ButtonAdd',
  props: [],
  data () {
    return {
    }
  },
  components: {},
  methods: {
    emitTap () {
      this.$emit('on-tap')
    }
  }
}
</script>

<template>
    <div class="button"  @click="emitTap">
        <span class="add-icon">
            <Icon type="plus"></Icon>
        </span>
        <span class="button-content">
            <slot></slot>
        </span>
    </div>
</template>

<style lang="less" scoped>
    .button {
        vertical-align: middle;
        font-size: 12px;

        height: 30px;
        line-height: 30px;
        border: 1px solid #D83F3F;
        border-radius: 100px;
        background: #ffffff;
        color: #D83F3F;

        display: inline-block;

        margin: 15px 10px 0 0;

        padding: 0 10px;

        cursor: pointer;

        user-select: none;
    }

</style>
