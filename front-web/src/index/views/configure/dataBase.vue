<template>
	<!-- eslint-disable -->
	<div>
		<div class="margin-b20">
			<Button class="margin-r10" @click="modal1 = true">增加</Button>
			<!--<Button  class="margin-r10"  @click="del">删除</Button>-->
			<Button class="margin-r10" @click="refresh">刷新表格</Button>
		</div>
		<f-table :data="newHostGroup" @on-selection-change="checkSelect">
			<!--<f-table-column type="selection" :width="60"></f-table-column>-->
			<f-table-column title="name" prop="name"></f-table-column>
			<f-table-column title="max_conn" prop="max_conn"></f-table-column>
			<f-table-column title="init_conn" prop="init_conn"></f-table-column>
			<f-table-column title="user" prop="user"></f-table-column>
			<f-table-column title="password" prop="password"></f-table-column>
			<f-table-column title="write" prop="write"></f-table-column>
			<f-table-column title="active_write" prop="active_write"></f-table-column>
			<f-table-column title="enable_switch" prop="enable_switch"></f-table-column>
			<f-table-column title="操作">
				<template slot-scope="scope">
					<inline-button @on-tap="del(scope)">删除</inline-button>
				</template>
			</f-table-column>
		</f-table>
		<Page class="page-bar" :total="tabelList.length" :page-size="10" @on-change="pageOnChange" show-total></Page>

		<Modal v-model="modal1" title="新增hostgroup配置" :closable="false">
			<Form
				ref="dataformValidate"
				:model="formValidate"
				label-position="left"
				:label-width="110"
				:rules="dataRuleValidate"
			>
				<FormItem label="name" prop="name"> <Input v-model="formValidate.name" /> </FormItem>
				<FormItem label="max_conn" prop="maxConn">
					<InputNumber v-model="formValidate.maxConn"></InputNumber>
				</FormItem>
				<FormItem label="init_conn" prop="initConn">
					<InputNumber v-model="formValidate.initConn"></InputNumber>
				</FormItem>
				<FormItem label="user" prop="user"> <Input v-model="formValidate.user" /> </FormItem>
				<FormItem label="password" prop="password">
					<Input type="password" v-model="formValidate.password" />
				</FormItem>
				<FormItem label="write" prop="write"> <Input v-model="formValidate.write" /> </FormItem>
				<FormItem label="read" prop="read"> <Input v-model="formValidate.read" /> </FormItem>
				<FormItem label="active_write" prop="activeWrite">
					<InputNumber v-model="formValidate.activeWrite"></InputNumber>
				</FormItem>
				<FormItem label="enable_switch" prop="enableSwitch">
					<i-switch v-model="formValidate.enableSwitch" />
				</FormItem>
			</Form>
			<div slot="footer">
				<Button type="text" size="large" @click="dataAddcancel">取消</Button>
				<Button type="primary" size="large" @click="dataAddok">确定</Button>
			</div>
		</Modal>
	</div>
</template>
<script>
import fTable from '@components/fTable/fTable'
import fTableColumn from '@components/fTable/fTableColumn'
import Crumb from '@components/crumb'
import { Mymixin } from '@utils/utils'
import InlineButton from '@components/button/inline-button'
import API from '@api'

export default {
	mixins: [Mymixin],
	data() {
		const validateInitConn = (rule, value, callback) => {
			if (typeof value === 'number' && (value > 0 && value <= this.formValidate.maxConn)) {
				callback()
			} else {
				callback(new Error('不能为空且大于0小于maxConn的数值'))
			}
		}
		const validateRead = (rule, value, callback) => {
			if (this.formValidate.write !== '') {
				let len = 0
				if (this.formValidate.write.indexOf(',') !== -1) {
					len = this.formValidate.write.split(',').length
				} else {
					len = 1
				}
				if (value >= 0 && value <= len - 1) {
					callback()
				} else {
					callback(new Error('不能为空且大于等于0小于write减一'))
				}
			} else {
				callback(new Error('write不能为空'))
			}
		}
		return {
			dataList: {},
			modal1: false,
			names: [],
			formValidate: {
				name: '',
				maxConn: 1024,
				initConn: 10,
				user: '',
				password: '',
				write: '',
				read: '',
				activeWrite: 0,
				enableSwitch: false
			},
			dataRuleValidate: {
				name: [{ required: true, message: '不能为空', trigger: 'blur' }],
				maxConn: [{ type: 'number', min: 1024, required: true, message: '不能为空且大于等于1024', trigger: 'blur' }],
				initConn: [{ required: true, validator: validateInitConn, trigger: 'blur' }],
				user: [{ required: true, message: '不能为空', trigger: 'blur' }],
				password: [{ required: true, message: '不能为空', trigger: 'blur' }],
				write: [{ required: true, message: '不能为空', trigger: 'blur' }],
				activeWrite: [{ required: true, validator: validateRead, trigger: 'blur' }],
				enableSwitch: [{ required: true }]
			}
		}
	},
	props: {
		tabelList: {
			type: Array,
			default() {
				return []
			}
		},
		newHostGroup: {
			type: Array,
			default() {
				return []
			}
		},
		clusterName: {
			type: String,
			default() {
				return ''
			}
		}
	},
	methods: {
		refresh() {
			this.$emit('resetTable')
		},
		checkSelect() {
			// this.names=[]
			// var that=this;
			//
			// for(var i=0;i<selection.length;i++){
			// 	that.names.push(selection[i].name)
			// }
		},
		del(scope) {
			const self = this
			self.names = []
			self.names.push(scope.name)
			if (self.names.length > 0) {
				API.delHostgroups({ clusterName: self.clusterName, names: self.names }).then(res => {
					if (res == null) {
						this.$Modal.success({
							title: '消息',
							content: '删除成功!'
						})
					} else {
						const d = res
						let content = ''
						for (let i = 0; i < d.length; i++) {
							content += `<p>${d[i]}</p>`
						}
						this.$Modal.success({
							title: '消息',
							content
						})
					}
					this.$emit('resetTable')
				})
			} else {
				self.$Message.info('必须添加大于一条的数据!')
			}
		},
		pageOnChange(page) {
			this.$emit('setTableList', page - 1)
		},
		dataAddok() {
			const that = this
			that.$refs.dataformValidate.validate(valid => {
				if (valid) {
					const data = {
						name: that.formValidate.name,
						max_conn: that.formValidate.maxConn,
						init_conn: that.formValidate.initConn,
						user: that.formValidate.user,
						password: that.formValidate.password,
						write: that.formValidate.write,
						read: that.formValidate.read,
						active_write: that.formValidate.activeWrite,
						enable_switch: that.formValidate.enableSwitch
					}

					const d = {}
					d.clusterName = that.clusterName
					d.groups = []
					d.groups.push(data)
					API.addHostgroups(d)
						.then(res => {
							if (res == null) {
								that.$Modal.success({
									title: '消息',
									content: '增加成功!'
								})
							} else {
								let content = ''
								for (let i = 0; i < res.length; i++) {
									content += `<p>${res[i]}</p>`
								}
								that.$Modal.success({
									title: '消息',
									content
								})
							}
							that.$emit('resetTable')
						})
						.catch(() => {
							that.$refs.dataformValidate.resetFields()
						})
					that.modal1 = false
				} else {
					that.modal1 = true
				}
			})
		},
		dataAddcancel() {
			this.$refs.dataformValidate.resetFields()
			this.modal1 = false
		}
	},
	components: {
		fTable,
		fTableColumn,
		Crumb,
		InlineButton
	}
}
</script>

<style scoped lang="less">
.configure-head {
	padding: 30px 30px 0px 30px;
	.padding15 {
		padding: 15px;
		background: #ffffff;
		-webkit-box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.13);
		box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.13);
		margin-bottom: 20px;
	}
}
.configure {
	background-color: #ffffff;
}
.container {
	padding: 30px;
}
.colony-name {
	margin-bottom: 10px;
	padding: 8px 16px;
	font-size: 16px;
}
</style>
