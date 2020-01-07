<template>
	<!-- eslint-disable -->
	<div>
		<div class="margin-b20">
			<Button class="margin-r10" @click="addSubTable = true">增加</Button>
			<!--<Button class="margin-r10" @click="delSchema">删除</Button>-->
			<Button class="margin-r10" @click="updateSchema" style="width: 100px">更新读写分离</Button>
			<Button  class="margin-r10"  @click="refresh">刷新表格</Button>
		</div>
		<f-table :data="newHostGroup" @on-selection-change="checkSelect">
			<f-table-column type="selection" :width="60"></f-table-column>
			<f-table-column title="name" prop="Name"></f-table-column>
			<f-table-column title="host_group_cluster" prop="hostgroupcluster"></f-table-column>
			<f-table-column title="rw_split" prop="rw_split"></f-table-column>
			<f-table-column title="multi_route" prop="multi_route"></f-table-column>
			<f-table-column title="init_conn_multi_route" prop="init_conn_multi_route"></f-table-column>
			<f-table-column title="schema_sharding" prop="schema_sharding"></f-table-column>
			<f-table-column title="table_sharding" prop="table_sharding"></f-table-column>
			<f-table-column title="table_configs">
				<template slot-scope="scope">
					<inline-button @on-tap='detail(scope)'>查看</inline-button>
					<inline-button @on-tap='delSchema(scope)'>删除</inline-button>
				</template>
			</f-table-column>
		</f-table>
		<Page class="page-bar" :total="tabelList.length"  :page-size='10' @on-change='pageOnChange'  show-total></Page>

		<Modal
			v-model="addSubTable"
			title="新增分库分表配置"
			:closable="false"
			width="700">
			<Form ref="subformValidate" :model="schemas[0]" label-position="left" :label-width="190" :rules="dataRuleValidate" >
				<FormItem label="name" prop="name">
					<Input v-model="schemas[0].name" placeholder="name" />
				</FormItem>
				<FormItem label="custody" prop="custody">
					<i-switch v-model="schemas[0].custody" />
				</FormItem>
				<FormItem label="HostGroupCluster" prop="hostGroupCluster">
					<Select v-model="hostGroupCluster" class="margin-b20" @on-change="hostGroupClusterChange" placeholder="群集名称">
						<Option v-for="(item,i) in hostGroupClusterList"  :value="item.name" :key="i">{{ item.name }}</Option>
					</Select>
					<Row :gutter="16" class="margin-b20" v-show="isTrue()">
						<Col span="10" class="text-r">nonsharding_host_group：</Col>
						<Col span="14">
						{{nonsharding_host_group}}
						</Col>
					</Row>
					<Row :gutter="16" class="margin-b20" v-show="isTrue()">
						<Col span="10" class="text-r">sharding_host_groups：</Col>
						<Col span="14">
						{{sharding_host_groups}}
						</Col>
					</Row>
				</FormItem>

				<FormItem label="schema_sharding" prop="schema_sharding">
					<InputNumber  v-model="schemas[0].schema_sharding"></InputNumber>
				</FormItem>
				<FormItem label="table_sharding" prop="table_sharding">
					<InputNumber  v-model="schemas[0].table_sharding"></InputNumber>
				</FormItem>
				<FormItem label="rw_split" prop="rw_split">
					<i-switch v-model="schemas[0].rw_split" />
				</FormItem>

				<FormItem label="table_configs" v-if="!schemas[0].custody">
					<FormItem style="margin-bottom: 20px;" label="name" prop="table_configs[0].name" :rules="{required: true, message: '不能为空', trigger: 'blur'}">
						<Input v-model="schemas[0].table_configs[0].name" />
					</FormItem>
					<FormItem style="margin-bottom: 20px;" label="sharding_key" prop="table_configs[0].sharding_key" :rules="{required: true, message: '不能为空', trigger: 'blur'}">
						<Input v-model="schemas[0].table_configs[0].sharding_key" />
					</FormItem>
					<FormItem style="margin-bottom: 20px;" label="rule" prop="table_configs[0].rule" :rules="{required: true, message: '不能为空', trigger: 'blur'}">
						<Select v-model="schemas[0].table_configs[0].rule">
							<Option value="int">int</Option>
							<Option value="string">string</Option>
						</Select>
					</FormItem>
				</FormItem>
			</Form>
			<div slot="footer">
				<Button type="text" size="large" @click="cancel">取消</Button>
				<Button type="primary" size="large" @click="addSchemaOk">确定</Button>
			</div>
		</Modal>
		<Modal
			v-model="detailModal"
			title="表配置"
			@on-ok="ok"
			width="800"
			@on-cancel="cancel">
			<Row :gutter="16" class="margin-b20">
				<Col span="3">
				<Button @click="addTableModel = true">增加</Button>
				</Col>
				<!--<Col span="3">-->
				<!--<Button @click="delTable">删除</Button>-->
				<!--</Col>-->
				<Col span="3">
				<Button @click="tableConfigRefresh">刷新表格</Button>
				</Col>
			</Row>
			<f-table :data="newtableConfigsList" @on-selection-change="configCheckSelect">
				<!--<f-table-column type="selection" :width="60"></f-table-column>-->
				<f-table-column title="name" prop="name"></f-table-column>
				<f-table-column title="sharding_key" prop="sharding_key"></f-table-column>
				<f-table-column title="rule" prop="rule"></f-table-column>
				<f-table-column title="table_configs">
					<template slot-scope="scope">
						<inline-button @on-tap='delTable(scope)'>删除</inline-button>
					</template>
				</f-table-column>
			</f-table>
			</f-table>
			<Page class="page-bar" v-if="hackReset" :total="tableConfigsList.length" :page-size='10' @on-change='settableConfigsList'  show-total></Page>
		</Modal>

		<Modal
			v-model="addTableModel"
			title="增加"
			:closable="false">
			<Form ref="tablelFormValidate" :model="tables[0]"  label-position="left"  :rules="ruleValidate" :label-width="100">
				<FormItem label="name" prop="name">
					<Input v-model="tables[0].name" ></Input>
				</FormItem>
				<FormItem label="sharding_key" prop="sharding_key">
					<Input v-model="tables[0].sharding_key"></Input>
				</FormItem>
				<FormItem label="rule" prop="rule">
					<Select v-model="tables[0].rule" >
						<Option value="int">int</Option>
						<Option value="string">string</Option>
					</Select>
				</FormItem>
			</Form>
			<div slot="footer">
				<Button type="text" size="large" @click="addTableCancel">取消</Button>
				<Button type="primary" size="large" @click="addTable">确定</Button>
			</div>
		</Modal>

		<Modal
			v-model="update"
			title="更新读写分离"
			@on-ok="updateShardrw">
			<Row :gutter="16" class="margin-b20" v-for="(val,key) in shardrw" :key="val+key">
				<Col span="4" class="text-r">{{key}}：</Col>
				<Col span="18">
				<i-switch v-model="shardrw[key]" />
				</Col>
			</Row>
		</Modal>
	</div>
</template>

<script>
import fTable from "@components/fTable/fTable"
import fTableColumn from "@components/fTable/fTableColumn"
import InlineButton from "@components/button/inline-button"
import Crumb from "@components/crumb"
import { Mymixin } from "@utils/utils"
import API from "../../../__BASE__/api"

export default {
    mixins: [Mymixin],
    data() {
        const IsPowerOfTwo = (rule, value, callback) => {
            /* eslint-disable */
            if (value > 0 && (value & (value - 1)) == 0) {
                callback()
            } else {
                callback(new Error("请填写2的n次幂"))
            }
        }
        return {
            hackReset: true,
            dataRuleValidate: {
                name: [{ required: true, message: "不能为空", trigger: "blur" }],
                custody: [{ required: true }],
                hostGroupCluster: [{ required: true, message: "不能为空", trigger: "blur" }],
                schema_sharding: [{ type: "number", validator: IsPowerOfTwo, required: true, trigger: "blur" }],
                table_sharding: [{ type: "number", validator: IsPowerOfTwo, required: true, trigger: "blur" }],
                rw_split: [{ required: true }]
            },
            ruleValidate: {
                name: [{ required: true, message: "不能为空", trigger: "blur" }],
                sharding_key: [{ required: true, message: "不能为空", trigger: "blur" }],
                rule: [{ required: true, message: "不能为空", trigger: "blur" }]
            },
            inputType: true,
            update: false,
            addSubTable: false,
            names: [],
            detailModal: false,
            tableConfigsList: [],
            newtableConfigsList: [],
            addTableModel: false,
            tables: [
                {
                    name: "",
                    sharding_key: "",
                    rule: ""
                }
            ],
            tableName: [],
            sharding_host_groups: "",
            nonsharding_host_group: "",
            schemas: [
                {
                    name: "",
                    custody: false,
                    hostGroupCluster: "",
                    schema_sharding: 128,
                    table_sharding: 1,
                    rw_split: false,
                    table_configs: [
                        {
                            name: "",
                            sharding_key: "",
                            rule: ""
                        }
                    ]
                }
            ],
            hostGroupCluster: "",
            Groups: [],
            shardrw: {},
            schemaName: "",
            index: 0
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
                return ""
            }
        },
        hostGroupClusterList: {
            type: Array,
            default() {
                return []
            }
        },
        snapshotName: {
            type: String,
            default() {
                return ""
            }
        }
    },
    methods: {
        isTrue() {
            if (typeof this.hostGroupCluster === "undefined") {
                return false
            } else {
                if (this.hostGroupCluster === "") {
                    return false
                } else {
                    return true
                }
            }
        },
        refresh() {
            this.$emit("resetTable")
        },
        tableConfigRefresh() {
            this.getConfigTable()
            this.hackReset = false
            this.$nextTick(() => {
                this.hackReset = true
            })
        },
        updateSchema() {
            let that = this
            if (this.names.length > 0) {
                that.update = true
            } else {
                this.$Message.info("必须选择不少于一条的数据!")
            }
        },
        updateShardrw() {
            API.updateShardrw({ clusterName: this.clusterName, shardrw: this.shardrw }).then(res => {
                if (res == null) {
                    this.$Modal.success({
                        title: "消息",
                        content: "更新成功!"
                    })
                } else {
                    var d = res
                    let content = ""
                    for (var i = 0; i < d.length; i++) {
                        content = content + `<p>${d[i]}</p>`
                    }
                    this.$Modal.success({
                        title: "消息",
                        content: content
                    })
                }
                this.$emit("resetTable")
            })
        },
        delSchema(scope) {
            this.names = []
            this.names.push(scope.Name)
            if (this.names.length > 0) {
                API.delSchema({ clusterName: this.clusterName, names: this.names }).then(res => {
                    if (res == null) {
                        this.$Modal.success({
                            title: "消息",
                            content: "删除成功!"
                        })
                    } else {
                        var d = res
                        let content = ""
                        for (var i = 0; i < d.length; i++) {
                            content = content + `<p>${d[i]}</p>`
                        }
                        this.$Modal.success({
                            title: "消息",
                            content: content
                        })
                    }
                    this.$emit("resetTable")
                })
            } else {
                this.$Message.info("必须添加大于一条的数据!")
            }
        },
        addSchemaOk() {
            let that = this
            this.$refs["subformValidate"].validate(valid => {
                if (valid) {
                    let data = {
                        clusterName: that.clusterName
                    }
                    if (that.schemas[0].custody) {
                        data.schemas = [
                            {
                                name: that.trim(that.schemas[0].name),
                                custody: that.schemas[0].custody,
                                hostGroupCluster: that.schemas[0].hostGroupCluster,
                                schema_sharding: that.schemas[0].schema_sharding,
                                table_sharding: that.schemas[0].table_sharding,
                                rw_split: that.schemas[0].rw_split
                            }
                        ]
                    } else {
                        data.schemas = [
                            {
                                name: that.trim(that.schemas[0].name),
                                custody: that.schemas[0].custody,
                                hostGroupCluster: that.schemas[0].hostGroupCluster,
                                schema_sharding: that.schemas[0].schema_sharding,
                                table_sharding: that.schemas[0].table_sharding,
                                rw_split: that.schemas[0].rw_split,
                                table_configs: [
                                    {
                                        name: that.schemas[0].table_configs[0].name,
                                        sharding_key: that.schemas[0].table_configs[0].sharding_key,
                                        rule: that.schemas[0].table_configs[0].rule
                                    }
                                ]
                            }
                        ]
                    }
                    API.addSchema(data).then(res => {
                        if (res == null) {
                            that.$Modal.success({
                                title: "消息",
                                content: "增加成功!"
                            })
                        } else {
                            var d = res
                            let content = ""
                            for (let i = 0; i < d.length; i++) {
                                content = content + `<p>${d[i]}</p>`
                            }
                            that.$Modal.success({
                                title: "消息",
                                content: content
                            })
                        }
                        that.$emit("resetTable")
                    })
                    that.$refs["subformValidate"].resetFields()
                    that.hostGroupCluster = ""
                    that.nonsharding_host_group = ""
                    that.sharding_host_groups = ""
                    that.addSubTable = false
                } else {
                    that.addSubTable = true
                }
            })
        },
        delTable(scope) {
            let self = this
            API.delTable({
                clusterName: self.clusterName,
                tableName: scope.name,
                schemaName: self.schemaName
            }).then(res => {
                if (res == null) {
                    self.$Modal.success({
                        title: "消息",
                        content: "删除成功!"
                    })
                } else {
                    let d = res
                    let content = ""
                    for (let i = 0; i < d.length; i++) {
                        content = content + `<p>${d[i]}</p>`
                    }
                    self.$Modal.success({
                        title: "消息",
                        content: content
                    })
                }

                self.getConfigTable()
                self.hackReset = false
                self.$nextTick(() => {
                    self.hackReset = true
                })
            })
        },
        addTable() {
            this.$refs["tablelFormValidate"].validate(valid => {
                if (valid) {
                    API.addTable({
                        clusterName: this.trim(this.clusterName),
                        schemaName: this.trim(this.schemaName),
                        tables: [
                            {
                                name: this.trim(this.tables[0].name),
                                sharding_key: this.trim(this.tables[0].sharding_key),
                                rule: this.trim(this.tables[0].rule)
                            }
                        ]
                    }).then(res => {
                        if (res == null) {
                            this.$Modal.success({
                                title: "消息",
                                content: "添加成功!"
                            })
                        } else {
                            var d = res
                            let content = ""
                            for (var i = 0; i < d.length; i++) {
                                content = content + `<p>${d[i]}</p>`
                            }
                            this.$Modal.success({
                                title: "消息",
                                content: content
                            })
                        }
                        this.getConfigTable()
                        this.hackReset = false
                        this.$nextTick(() => {
                            this.hackReset = true
                        })
                    })
                    this.addTableModel = false
                    this.$refs["tablelFormValidate"].resetFields()
                } else {
                    this.addTableModel = true
                }
            })
        },
        addTableCancel() {
            this.$refs["tablelFormValidate"].resetFields()
            this.addTableModel = false
        },
        settableConfigsList(page) {
            if (this.tableConfigsList.length > 0) {
                let p = page - 1
                this.newtableConfigsList = this.tableConfigsList.slice(p * 10, (p + 1) * 10)
            } else {
                this.newtableConfigsList = []
            }
        },
        detail(item) {
            this.index = item._index
            this.detailModal = true
            this.schemaName = item.Name
            this.getConfigTable()
        },
        getConfigTable() {
            API.getConfigTable({
                clusterName: this.clusterName,
                snapshotName: this.snapshotName,
                schemaName: this.schemaName
            }).then(res => {
                if (res == null) {
                    this.tableConfigsList = []
                } else {
                    this.tableConfigsList = res
                }
                this.settableConfigsList(1)
            })
        },
        configCheckSelect(selection) {
            this.tableName = []
            var that = this
            for (var i = 0; i < selection.length; i++) {
                that.tableName.push(selection[i].name)
            }
        },
        checkSelect(selection) {
            this.names = []
            var that = this
            for (var i = 0; i < selection.length; i++) {
                that.names.push(selection[i].Name)
                that.shardrw[selection[i].Name] = selection[i].rw_split
            }
        },
        pageOnChange(page) {
            this.$emit("setTableList", page - 1)
        },
        ok() {},
        cancel() {
            this.$refs["subformValidate"].resetFields()
            this.addSubTable = false
        },
        hostGroupClusterChange(value) {
            let that = this
            if (value !== "") {
                let hgcl = this.hostGroupClusterList
                for (let i = 0; i < hgcl.length; i++) {
                    if (hgcl[i].name === value) {
                        that.sharding_host_groups = hgcl[i].sharding_host_groups
                        that.nonsharding_host_group = hgcl[i].nonsharding_host_group
                        that.schemas[0].hostGroupCluster = hgcl[i].name
                    }
                }
            }
        }
    },
    components: {
        fTable,
        fTableColumn,
        Crumb,
        InlineButton
    },
    watch: {
        // tabelList (val,oldval){
        //
        // }
    }
}
</script>

<style scoped lang="less">
.group {
    display: inline-block;
    border: 1px solid #ddd;
    margin: 10px;
    padding: 10px;
    position: relative;
    .group-del {
        position: absolute;
        right: -5px;
        top: -5px;
    }
}
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
