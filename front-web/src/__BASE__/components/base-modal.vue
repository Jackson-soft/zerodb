<script>
/**
     -- zeqi@2dfire

     @props:
     show: 弹窗是否展示
     header: header显示文字内容

     @event:
     on-ok: 右侧确定按钮 点击事件
     on-cancel: 左侧取消按钮 点击事件

     基础弹窗容器，只提供背景遮罩和白色内容框
     */

export default {
  name: 'baseModal',
  props: ['show', 'header', 'okText', 'cancelText', 'width'],
  data () {
    return {
    }
  },
  methods: {
    emitOk () {
      this.$emit('on-ok')
    },
    emitCancel () {
      this.$emit('on-cancel')
    }
  },
  computed: {
    cancelButton: function () {
      if (this.cancelText === undefined) {
        return '取消'
      } else {
        return this.cancelText
      }
    },
    okButton: function () {
      if (this.okText === undefined) {
        return '确定'
      } else {
        return this.okText
      }
    }
  },
  created () {

  }
}
</script>

<template>
    <div class="base-modal">
        <Modal :value="show" :closable="false" :mask-closable="false" class="modal" class-name="vertical-center-modal" :width="width">
            <p slot="header">
                {{header}}
            </p>
            <slot>
                <!--内部插槽-->
            </slot>
            <div slot="footer">
                <Button type="ghost" class="modal-btn" @click.native="emitCancel">{{cancelButton}}</Button>
                <Button type="primary" class="modal-btn" @click.native="emitOk">{{okButton}}</Button>
            </div>
        </Modal>
    </div>
</template>

<style lang="less" scoped>
    .modal {
        font-size: 12px;
    }
     .modal /deep/ .ivu-modal-body{
        min-height: 300px;
    }
    .modal-btn {
        width: 88px;
    }

</style>
