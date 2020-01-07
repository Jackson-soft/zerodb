<script>
/**
     -- zipai@2dfire

     @props:

     @events:
     on-tap: 点击的冒泡事件

     @disabled:
     是否禁用

     */

export default {
	name: 'InlineButton',
	props: ['isdisable', 'colorType'],
	data() {
		return {}
	},
	components: {},
	methods: {
		emitTap() {
			if (!this.isdisable) {
				this.$emit('on-tap')
			}
		}
	}
}
</script>

<template>
	<!-- eslint-disable -->
	<div class="button">
		<button :class="['inline-buttom', { disabled: this.isdisable, colorType: this.colorType }]" @click="emitTap">
			<slot></slot>
		</button>
	</div>
</template>

<style lang="less" scoped>
.button {
	display: inline-block;
}
.inline-buttom {
	border: none;
	cursor: pointer;
	outline: none;
	color: #5599ff;
	background: transparent;
}
.colorType {
	color: #cccccc;
}
.disabled {
	color: #cccccc;
	cursor: not-allowed;
}
</style>
