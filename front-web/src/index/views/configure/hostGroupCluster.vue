<template>
	<!-- eslint-disable -->
	<div>
		<div class="margin-b20">
			<Button class="margin-r10" @click="modal1 = true">增加</Button>
			<!--<Button class="margin-r10" @click="hostGroupdel">删除</Button>-->
			<Button  class="margin-r10"  @click="refresh">刷新表格</Button>
		</div>
		<f-table :data="newHostGroup" @on-selection-change="checkSelect">
			<!--<f-table-column type="selection" :width="60"></f-table-column>-->
			<f-table-column title="name" prop="name"></f-table-column>
			<f-table-column title="nonsharding_host_group" prop="nonsharding_host_group"></f-table-column>
			<f-table-column title="sharding_host_groups" prop="sharding_host_groups"></f-table-column>
			<f-table-column title="操作" >
				<template slot-scope="scope">
					<inline-button @on-tap='updateHostGroupCluster(scope)'>更新</inline-button>
					<inline-button @on-tap='hostGroupdel(scope)'>删除</inline-button>
				</template>
			</f-table-column>
		</f-table>
		<Page class="page-bar" :total="tabelList.length" :page-size='10' @on-change='pageOnChange'  show-total></Page>


		<Modal
			v-model="modal1"
			:closable="false"
			title="新增HostGroupCluster配置">
			<Form ref="clusterformValidate" :model="formValidate" label-position="left" :label-width="190" :rules="dataRuleValidate" >
				<FormItem label="clusterName" prop="clusterName">
					{{clusterName}}
				</FormItem>
				<FormItem label="HostGroupClusterName" prop="HostGroupClusterName">
					<Input v-model="formValidate.HostGroupClusterName"/>
				</FormItem>
				<FormItem label="nonsharding_Host_Group" prop="nonsharding_Host_Group">
					<Select v-model="formValidate.nonsharding_Host_Group" class="margin-b20" clearable placeholder="nonsharding_Host_Group">
						<Option v-for="(item,i) in hostGroupList"  :value="item.name" :key="i">{{ item.name }}</Option>
					</Select>
				</FormItem>
				<FormItem label="hostGroup" prop="hostGroup">
					<Select v-model="formValidate.hostGroup" @on-change="hostGroupChange" clearable class="margin-b20" placeholder="群集名称">
						<Option v-for="(item,i) in hostGroupList"  :value="item.name" :key="i">{{ item.name }}</Option>
					</Select>
					<div v-for="(item,i) in Groups" class="group" :key="i">
						<Icon class="group-del" @click="del(i)" type="ios-close" size="14" />
						<div>{{item}}</div>
					</div>
				</FormItem>
			</Form>
			<div slot="footer">
				<Button type="text" size="large" @click="dataAddcancel">取消</Button>
				<Button type="primary" size="large" @click="dataAddok">确定</Button>
			</div>
		</Modal>

		<Modal
			v-model="update"
			title="更新HostGroupCluster配置"
			@on-ok="updateok">
			<Row :gutter="16" class="margin-b20" v-for="(val, key, index) in HostGroupCluster" :key="key">
				<Col span="8" class="text-r">{{key}}：</Col>
				<Col span="14">
				<Input v-if="key=='name'" v-model="HostGroupCluster[key]" />
				<Select v-else-if="key=='nonsharding_host_group'" v-model="HostGroupCluster[key]" clearable class="margin-b20" placeholder="nonsharding_Host_Group">
					<Option v-for="(item,i) in hostGroupList"  :value="item.name" :key="i">{{ item.name }}</Option>
				</Select>
				<div v-else>
					<Select @on-change="changeCluster" class="margin-b20" clearable placeholder="群集名称">
						<Option v-for="(item,i) in hostGroupList"  :value="item.name" :key="i">{{ item.name }}</Option>
					</Select>
					<div v-for="(item,i) in HostGroupCluster[key]" class="group" :key="i">
						<Icon class="group-del" @click="delUpdate(i)" type="ios-close" size="14" />
						<div>{{item}}</div>
					</div>
				</div>
				</Col>
			</Row>
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
		return {
			dataList: {},
			names: [],
			modal1: false,
			Groups: [],
			HostGroupCluster: {},
			update: false,
			formValidate: {
				HostGroupClusterName: '',
				nonsharding_Host_Group: '',
				hostGroup: '',
			},
			dataRuleValidate: {
				HostGroupClusterName: [
					{ required: true, message: '不能为空', trigger: 'blur' }
				],
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
		hostGroupList: {
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
		},
		schemaName: {
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
		hostGroupChange(value) {
			const self = this
			if (typeof(value) !== 'undefined') {
				if (self.Groups.indexOf(value) === -1) {
					self.Groups.push(value)
				}
			}
		},
		changeCluster(value) {
			if (this.HostGroupCluster.sharding_host_groups.indexOf(value) === -1) {
				this.HostGroupCluster.sharding_host_groups.push(value)
			}
		},
		updateok() {
			const data = {}
			const that = this
			data.ClusterName = that.clusterName
			data.HostGroupCluster = that.HostGroupCluster
			API.updateHostGroupCluster(data).then((res) => {
				if (res == null) {
					that.$Modal.success({
						title: '消息',
						content: '更新成功!',
					})
				} else {
					const d = res
					let content = ''
					for (let i = 0; i < d.length; i++) {
						content += `<p>${d[i]}</p>`
					}
					that.$Modal.success({
						title: '消息',
						content
					})
				}
				that.$emit('resetTable')
			})
		},
		updateHostGroupCluster(item) {
			this.update = true
			this.HostGroupCluster = {
				name: item.name,
				nonsharding_host_group: item.nonsharding_host_group,
				sharding_host_groups: item.sharding_host_groups
			}
		},
		checkSelect(selection) {
			const that = this
			that.names = []
			for (let i = 0; i < selection.length; i++) {
				that.names.push(selection[i].name)
			}
		},
		pageOnChange(page) {
			this.$emit('setTableList', page - 1)
		},
		dataAddok() {
			const that = this
			that.$refs.clusterformValidate.validate((valid) => {
				if (valid) {
					API.addHostgroupcluster({
						ClusterName: that.trim(that.clusterName),
						HostGroupCluster: {
							Name: that.trim(that.formValidate.HostGroupClusterName),
							nonsharding_Host_Group: that.formValidate.nonsharding_Host_Group,
							sharding_Host_Groups: that.Groups
						}
					}).then((res) => {
						if (res == null) {
							that.$Modal.success({
								title: '消息',
								content: '增加成功!',
							})
						} else {
							const d = res
							let content = ''
							for (let i = 0; i < d.length; i++) {
								content += `<p>${d[i]}</p>`
							}
							that.$Modal.success({
								title: '消息',
								content
							})
						}
						that.$emit('resetTable')
						that.dataAddcancel()
					}).catch(() => {
						that.dataAddcancel()
					})
				} else {
					that.modal1 = true
				}
			})
		},
		dataAddcancel() {
			this.$refs.clusterformValidate.resetFields()
			this.Groups = []
			this.modal1 = false
		},
		del(index) {
			this.Groups.splice(index, 1)
		},
		delUpdate(index) {
			this.HostGroupCluster.sharding_host_groups.splice(index, 1)
		},
		hostGroupdel(scope) {
			const that = this
			that.names = []
			that.names.push(scope.name)
			/* eslint-disable */
			if(that.names.length>0){
				API.delHostgroupcluster({ClusterName:that.clusterName,Snapshot:that.schemaName,HostGroupClusterName:that.names}).then((res)=>{
					if(res==null){
						that.$Modal.success({
							title: '消息',
							content: '删除成功!',
						});
					}else{
						let d=res;
						let content= '';
						for (let i=0;i<d.length;i++){
							content=content+`<p>${d[i]}</p>`
						}
						that.$Modal.success({
							title: '消息',
							content: content
						});
					}
					that.$emit('resetTable')
				})
			}else{
				that.$Message.info('必须添加大于一条的数据!');
			}
		}
	},
	components: {
		fTable,
		fTableColumn,
		Crumb,
		InlineButton
	},
	watch:{
		// formValidate: {
		// 	handler(newName, oldName) {
		//
		// 	},
		// 	deep: true
		// }
	}
}
</script>

<style scoped lang="less">
	.group{
		display: inline-block;
		border:1px solid #ddd;
		margin: 10px;
		padding:10px;
		position: relative;
		.group-del{
			position: absolute;
			right: -5px;
			top:-5px;
		}
	}
	.configure-head{
		padding:30px 30px 0px 30px;
		.padding15{
			padding: 15px;
			background: #ffffff;
			-webkit-box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.13);
			box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.13);
			margin-bottom: 20px;
		}
	}
	.configure{
		background-color: #ffffff;
	}
	.container{
		padding: 30px;
	}
	.colony-name{
		margin-bottom: 10px;
		padding:8px 16px;
		font-size: 16px;
	}

</style>
