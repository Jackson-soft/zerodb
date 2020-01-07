<script>
/**
     标签组件 zipai at 2018/5/11

     传入标签内容 text（左侧显示，必填）和 type: 1 显示必选小红点 ，默认不显示必选小红点

     */

let componentData = {
  text: 'label',
  type: 0,
  size: 'md'
}

export default {
  data () {
    return {
      componentData: componentData
    }
  },
  props: ['text', 'type', 'size'],
  methods: {},
  created () {

  }
}
</script>
<template>
    <div class="open-label md" :class="{'lg':size=='lg'}">
        <span class="must" v-if="type==1"></span>
        <span class="needless" v-if="type==0"></span>
        <span class="label-text">{{text}}</span>
    </div>
</template>

<style lang="less" scoped>

    .open-label {
        text-align: right;
        display: inline-block;
        vertical-align: top;
    }

    .md {
        font-size: 13px;
        color: #333333;
        line-height: 12px;
        margin-top: 10px;
        .label-text{
            height: 12px;
        }
    }
    .lg {
        font-size: 16px;
        line-height: 16px;
        margin-top: 8px;
        .label-text{
            height: 16px;
        }
    }

    .must {
        width: 4px;
        height: 4px;
        background: #B11A1B;
        border-radius: 50%;
        display: inline-block;
        vertical-align: middle;
    }
     .needless {
        width: 4px;
        height: 4px;
        background:transparent;
        border-radius: 50%;
        display: inline-block;
        vertical-align: middle;
    }
</style>
