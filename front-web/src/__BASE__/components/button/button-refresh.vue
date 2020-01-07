<script>
/**
     -- zeqi@2dfire

     @props:

     @events:
     on-tap: 点击的冒泡事件

     @desc:
     刷新按钮

     */

export default {
  name: 'ButtonRefresh',
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
        <span class="refresh">
        </span>
    </div>
</template>

<style lang="less" scoped>
    .button {
        width: 36px;
        height: 36px;

        background: #FFFFFF;
        border: 1px solid #F1F1F1;
        box-shadow: 0 2px 2px 0 rgba(0,0,0,0.10);

        border-radius: 50%;

        cursor: pointer;

        user-select: none;

        text-align: center;
    }

    .refresh{
        display: inline-block;
        height: 36px;
        width: 20px;
        background: url("https://assets.2dfire.com/frontend/18da573dfbeb468490ab620159d90c61.png");

        background-size: 20px 20px;
        background-repeat: no-repeat;
        background-position-y: 50%;

    }

</style>
