<template>
    <!-- eslint-disable -->
    <div>
        <Crumb></Crumb>
        <div class="configure-head">
            <div class="padding15">
                <Row>
                    <Col :sm="8" :md="10" :lg="10">
                        <Select
                            class="margin-r10 configure-w"
                            v-model="clusterName"
                            @on-change="clusterNameChage(clusterName)"
                            @on-open-change="getClusterName"
                            placeholder="群集名称"
                        >
                            <Option v-for="(item, i) in clusterList" :value="item" :key="item + i">{{ item }}</Option>
                        </Select>
                        <!--<Select class="configure-w" v-model="snapshotName" @on-change="snapshotNameChage" @on-open-change="getSnapshotName" placeholder="快照名称" >-->
                        <!--<Option v-for="list in snapshotList" :value="list" :key="list">{{ list }}</Option>-->
                        <!--</Select>-->
                    </Col>
                    <Col :sm="16" :md="14" :lg="14">
                        <Button class="margin-r10" @click="snapshotModal = true">创建快照</Button>
                        <Button class="margin-r10" @click="configureModal = true">推送配置</Button>
                        <Button class="margin-r10" @click="initModal = true">初始化配置</Button>
                        <Button class="margin-r10" @click="exportModal = true">导出</Button>
                        <Button class="margin-r10" @click="rollbackConfig">回滚配置</Button>
                    </Col>
                </Row>
            </div>
        </div>
        <div class="container">
            <padding-view>
                <div>
                    <div class="colony-name">{{ clusterName ? clusterName : "集群名称" }}</div>
                    <Tabs v-model="newTabName" :animated="false" @on-click="tabClick">
                        <TabPane label="基本信息" name="name1">
                            <div class="margin-t20">
                                <Row
                                    :gutter="16"
                                    class="margin-b20"
                                    v-for="(val, key, index) in dataList"
                                    :key="val + key"
                                >
                                    <Col span="6" class="text-r">{{ key }}：</Col>
                                    <Col span="18" v-if="key == 'password'">******</Col>
                                    <Col span="18" v-else>{{ val }}</Col>
                                </Row>
                            </div>
                            <Button v-show="JSON.stringify(dataList) != '{}'" @click="updateBasic = true">更新</Button>
                        </TabPane>
                        <TabPane label="停服配置" name="name2">
                            <div class="margin-t20">
                                <Row
                                    :gutter="16"
                                    class="margin-b20"
                                    v-for="(val, key, index) in stopServiceList"
                                    :key="val + key"
                                >
                                    <Col span="6" class="text-r">{{ key }}：</Col>
                                    <Col span="18">{{ val }}</Col>
                                </Row>
                            </div>
                            <Button v-show="JSON.stringify(stopServiceList) != '{}'" @click="updateStopservice = true"
                                >更新</Button
                            >
                        </TabPane>
                        <TabPane label="切换配置" name="name3">
                            <div class="margin-t20">
                                <Row
                                    :gutter="16"
                                    class="margin-b20"
                                    v-for="(val, key, index) in switchConfigList"
                                    :key="val + key"
                                >
                                    <Col span="6" class="text-r">{{ key }}：</Col>
                                    <Col span="18">{{ val }}</Col>
                                </Row>
                            </div>
                            <Button v-show="JSON.stringify(switchConfigList) != '{}'" @click="updateSwitch = true"
                                >更新</Button
                            >
                        </TabPane>
                        <TabPane label="HostGroup配置" name="name4">
                            <data-base
                                :tabelList="hostGroupList"
                                :newHostGroup="newHostGroup"
                                :clusterName="clusterName"
                                @setTableList="sethostGroupList"
                                @resetTable="hostgroups"
                            ></data-base>
                        </TabPane>
                        <TabPane label="HostGroupCluster配置" name="name5">
                            <host-group-cluster
                                :hostGroupList="hostGroupList"
                                :tabelList="hostgroupClustersList"
                                :newHostGroup="newClustersList"
                                :clusterName="clusterName"
                                :schemaName="snapshotName"
                                @setTableList="setClustersList"
                                @resetTable="hostgroupClusters"
                            ></host-group-cluster>
                        </TabPane>
                        <TabPane label="分库分表配置" name="name6">
                            <sub-table
                                :hostGroupClusterList="hostgroupClustersList"
                                :tabelList="pullList"
                                :newHostGroup="newPullList"
                                :snapshotName="snapshotName"
                                :clusterName="clusterName"
                                @setTableList="setPullList"
                                @resetTable="full"
                            ></sub-table>
                        </TabPane>
                    </Tabs>
                </div>
            </padding-view>
        </div>
        <Modal v-model="updateBasic" title="基本信息更新">
            <Form
                ref="formValidate"
                label-position="left"
                :model="dataList"
                :label-width="110"
                :rules="baseRuleValidate"
            >
                <FormItem :label="key" :prop="key" v-for="(val, key, index) in dataList" :key="key">
                    <i-switch v-if="typeof val == 'boolean'" v-model="dataList[key]" />
                    <Input v-else-if="key == 'password'" type="password" v-model="dataList[key]" />
                    <InputNumber v-else-if="key == 'slow_log_time'" v-model="dataList[key]"></InputNumber>
                    <Input v-else v-model="dataList[key]" />
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="updateBasicCancel">取消</Button>
                <Button type="primary" size="large" @click="updateBasicOk">确定</Button>
            </div>
        </Modal>
        <Modal v-model="updateStopservice" title="停服配置更新">
            <Form
                ref="updateStopserviceRef"
                label-position="left"
                :model="stopServiceList"
                :label-width="190"
                :rules="stopRuleValidate"
            >
                <FormItem :label="key" :prop="key" v-for="(val, key, index) in stopServiceList" :key="key">
                    <i-switch v-if="typeof val == 'boolean'" v-model="stopServiceList[key]" />
                    <InputNumber
                        v-else-if="key == 'offline_swh_rejected_num'"
                        v-model="stopServiceList[key]"
                    ></InputNumber>
                    <InputNumber
                        v-else-if="key == 'offline_down_host_num'"
                        v-model="stopServiceList[key]"
                    ></InputNumber>
                    <Input v-else v-model="stopServiceList[key]" />
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="updateStopserviceCancel">取消</Button>
                <Button type="primary" size="large" @click="updateStopserviceOk">确定</Button>
            </div>
        </Modal>
        <Modal v-model="updateSwitch" title="切换配置更新">
            <Form
                ref="switchStopserviceRef"
                label-position="left"
                :model="switchConfigList"
                :label-width="190"
                :rules="switchRuleValidate"
            >
                <FormItem :label="key" :prop="key" v-for="(val, key, index) in switchConfigList" :key="key">
                    <i-switch v-if="typeof val == 'boolean'" v-model="switchConfigList[key]" />
                    <InputNumber v-else v-model="switchConfigList[key]"></InputNumber>
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="updateSwitchCancel">取消</Button>
                <Button type="primary" size="large" @click="updateSwitchOk">确定</Button>
            </div>
        </Modal>
        <Modal v-model="snapshotModal" title="创建快照" :closable="false">
            <Form ref="formSnapshotValidate" :model="snapshotValidate" :label-width="80">
                <FormItem
                    label="集群名称"
                    prop="addClusterName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select v-model="snapshotValidate.addClusterName" placeholder="群集名称">
                        <Option v-for="(item, i) in clusterList" :value="item" :key="item + i">{{ item }}</Option>
                    </Select>
                </FormItem>
                <FormItem
                    label="快照名称"
                    prop="addSnapshotName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Input
                        v-model="snapshotValidate.addSnapshotName"
                        placeholder="默认生成快照名(集群名_yyyy-dd-mm-HH--mm-ss)"
                    />
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="creatSnapshotCancel">取消</Button>
                <Button type="primary" size="large" @click="creatSnapshotOk">确定</Button>
            </div>
        </Modal>

        <Modal v-model="configureModal" title="推送配置" :closable="false">
            <Form ref="formConfigureValidate" :model="configureValidate" :label-width="80">
                <FormItem
                    label="集群名称"
                    prop="configureClusterName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select
                        v-model="configureValidate.configureClusterName"
                        @on-change="getSnapshotList(configureValidate.configureClusterName)"
                        placeholder="群集名称"
                    >
                        <Option v-for="(item, i) in clusterList" :value="item" :key="item + i">{{ item }}</Option>
                    </Select>
                </FormItem>
                <FormItem
                    label="快照名称"
                    prop="configureSnapshotName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select v-model="configureValidate.configureSnapshotName" placeholder="快照名称">
                        <Option v-for="item in snapshotList" :value="item" :key="item">{{ item }}</Option>
                    </Select>
                </FormItem>
                <FormItem label="快照名称" prop="pushTime">
                    <DatePicker
                        v-model="configureValidate.pushTime"
                        :options="options1"
                        @on-change="pushTimeChange"
                        format="yyyy-MM-dd HH:mm:ss"
                        type="datetime"
                        placeholder="Select date and time"
                    ></DatePicker>
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="configureCancel">取消</Button>
                <Button type="primary" size="large" @click="configureOk">确定</Button>
            </div>
        </Modal>
        <Modal v-model="initModal" title="初始化配置" :closable="false">
            <Form ref="formInitValidate" :model="initValidate" :label-width="80">
                <FormItem
                    label="集群名称"
                    prop="initClusterName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Input v-model="initValidate.initClusterName" />
                </FormItem>
                <FormItem label="上传文件" prop="file">
                    <input type="file" ref="inputer" @change="uploading()" class="file" />
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="initCancel">取消</Button>
                <Button type="primary" size="large" @click="initOk">确定</Button>
            </div>
        </Modal>

        <Modal v-model="exportModal" title="导出配置文件" :closable="false">
            <Form ref="formExportValidate" :model="exportValidate" :label-width="80">
                <FormItem
                    label="集群名称"
                    prop="exportClusterName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select
                        v-model="exportValidate.exportClusterName"
                        @on-change="getSnapshotList(exportValidate.exportClusterName)"
                        placeholder="群集名称"
                    >
                        <Option v-for="(item, i) in clusterList" :value="item" :key="item + i">{{ item }}</Option>
                    </Select>
                </FormItem>
                <FormItem label="快照名称" prop="exportSnapshotName">
                    <Select v-model="exportValidate.exportSnapshotName" placeholder="快照名称">
                        <Option v-for="item in snapshotList" :value="item" :key="item">{{ item }}</Option>
                    </Select>
                </FormItem>
                <FormItem
                    label="导出名称"
                    prop="exportName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Input v-model="exportValidate.exportName" placeholder="名称" />
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="exportCancel">取消</Button>
                <Button type="primary" size="large" @click="exportOk">确定</Button>
            </div>
        </Modal>
    </div>
</template>

<script>
import dataBase from "./dataBase"
import hostGroupCluster from "./hostGroupCluster"
import subTable from "./subTable"
import { mapGetters } from "vuex"
import Crumb from "@components/crumb"
import PaddingView from "@components/padding-view"
import { Mymixin } from "@utils/utils"
import API from "@api"

export default {
    mixins: [Mymixin],
    data() {
        return {
            snapshotValidate: {
                addClusterName: "",
                addSnapshotName: ""
            },
            configureValidate: {
                configureClusterName: "",
                configureSnapshotName: "",
                pushTime: ""
            },
            exportValidate: {
                exportClusterName: "",
                exportSnapshotName: "",
                exportName: ""
            },
            initValidate: {
                initClusterName: "",
                file: ""
            },
            baseRuleValidate: {
                config_name: [{ required: true, message: "config_name不能为空", trigger: "blur" }],
                user: [{ required: true, message: "user不能为空", trigger: "blur" }],
                password: [{ required: true, message: "password不能为空", trigger: "blur" }],
                slow_log_time: [
                    {
                        type: "number",
                        min: 1,
                        required: true,
                        message: "不能为空且必须是大于0的数值类型",
                        trigger: "blur"
                    }
                ],
                push_when_fail: [{ required: true }]
            },
            stopRuleValidate: {
                offline_on_lost_keeper: [{ required: true }],
                offline_swh_rejected_num: [
                    {
                        type: "number",
                        min: 1,
                        required: true,
                        message: "不能为空且必须是大于0的数值类型",
                        trigger: "blur"
                    }
                ],
                offline_down_host_num: [
                    {
                        type: "number",
                        min: 1,
                        required: true,
                        message: "不能为空且必须是大于0的数值类型",
                        trigger: "blur"
                    }
                ],
                offline_recover: [{ required: true }]
            },
            switchRuleValidate: {
                need_vote: [{ required: true }],
                vote_approve_ratio: [
                    {
                        type: "number",
                        min: 1,
                        max: 100,
                        required: true,
                        message: "不能为空且必须是大于0小于等于100的数值类型",
                        trigger: "blur"
                    }
                ],
                need_load_check: [{ required: true }],
                safe_load: [
                    {
                        type: "number",
                        min: 1,
                        required: true,
                        message: "不能为空且必须是大于0的数值类型",
                        trigger: "blur"
                    }
                ],
                need_binlog_check: [{ required: true }],
                safe_binlog_delay: [
                    {
                        type: "number",
                        min: 1,
                        required: true,
                        message: "不能为空且必须是大于0的数值类型",
                        trigger: "blur"
                    }
                ],
                binlog_wait_time: [
                    {
                        type: "number",
                        min: 1,
                        required: true,
                        message: "不能为空且必须是大于0的数值类型",
                        trigger: "blur"
                    }
                ],
                frequency: [
                    {
                        type: "number",
                        min: 1,
                        required: true,
                        message: "不能为空且必须是大于0的数值类型",
                        trigger: "blur"
                    }
                ],
                backend_ping_interval: [
                    {
                        type: "number",
                        min: 0,
                        required: true,
                        message: "不能为空且必须是大于等于0的数值类型",
                        trigger: "blur"
                    }
                ]
            },
            acceptType: ["yaml"],
            clusterName: "",
            snapshotName: "",
            dataList: {},
            snapshotModal: false,
            configureModal: false,
            initModal: false,
            exportModal: false,
            stopServiceList: {},
            newTabName: "name1",
            switchConfigList: {},
            // snapshotList:[]
            hostGroupList: [],
            newHostGroup: [],
            hostgroupClustersList: [],
            newClustersList: [],
            pullList: [],
            newPullList: [],
            updateStopservice: false,
            updateBasic: false,
            updateSwitch: false,
            options1: {
                disabledDate(date) {
                    return date && date.valueOf() < Date.now() - 86400000
                }
            },
            Basicloading: true
        }
    },
    methods: {
        getClusterName(val) {
            if (val === true) {
                this.$store.dispatch("zerodb/getClusterList")
            }
        },
        getSnapshotName(val) {
            if (val === true) {
                this.getSnapshotList(this.clusterName)
            }
        },
        rollbackConfig() {
            const self = this
            if (self.clusterName === "") {
                self.$Modal.warning({
                    title: "警告!",
                    content: "请选择集群！"
                })
            } else {
                API.rollbackConfig({ clusterName: self.clusterName }).then(res => {
                    if (res == null) {
                        self.$Modal.success({
                            title: "回滚成功!"
                        })
                    } else {
                        const d = res
                        let content = ""
                        for (let i = 0; i < d.length; i++) {
                            content += `<p>${d[i]}</p>`
                        }
                        self.$Modal.success({
                            title: "消息",
                            content
                        })
                    }
                })
            }
        },
        updateSwitchOk() {
            const that = this
            that.$refs.switchStopserviceRef.validate(valid => {
                if (valid) {
                    const data = {}
                    data.clusterName = that.clusterName
                    data.switch = that.switchConfigList
                    API.updateSwitch(data)
                        .then(res => {
                            if (res == null) {
                                that.$Modal.success({
                                    title: "消息",
                                    content: "操作成功!"
                                })
                            } else {
                                const d = res
                                let content = ""
                                for (let i = 0; i < d.length; i++) {
                                    content += `<p>${d[i]}</p>`
                                }
                                that.$Modal.success({
                                    title: "消息",
                                    content
                                })
                            }
                            that.interface()
                        })
                        .catch(() => {
                            that.interface()
                        })
                    that.updateSwitch = false
                } else {
                    that.updateSwitch = true
                }
            })
        },
        updateSwitchCancel() {
            this.$refs.switchStopserviceRef.resetFields()
            this.updateSwitch = false
        },
        updateBasicOk() {
            const that = this
            that.$refs.formValidate.validate(valid => {
                if (valid) {
                    const data = {}
                    data.clusterName = that.clusterName
                    data.basic = that.dataList
                    API.updateBasic(data)
                        .then(res => {
                            if (res == null) {
                                that.$Modal.success({
                                    title: "消息",
                                    content: "操作成功!"
                                })
                            } else {
                                const d = res
                                let content = ""
                                for (let i = 0; i < d.length; i++) {
                                    content += `<p>${d[i]}</p>`
                                }
                                that.$Modal.success({
                                    title: "消息",
                                    content
                                })
                            }
                            that.interface()
                        })
                        .catch(() => {
                            that.interface()
                        })
                    that.updateBasic = false
                } else {
                    that.updateBasic = true
                }
            })
        },
        updateBasicCancel() {
            this.$refs.formValidate.resetFields()
            this.updateBasic = false
        },
        updateStopserviceOk() {
            const that = this
            that.$refs.updateStopserviceRef.validate(valid => {
                if (valid) {
                    const data = {}
                    data.clusterName = that.clusterName
                    data.service = that.stopServiceList
                    API.updateStopservice(data)
                        .then(res => {
                            if (res == null) {
                                that.$Modal.success({
                                    title: "消息",
                                    content: "操作成功!"
                                })
                            } else {
                                const d = res
                                let content = ""
                                for (let i = 0; i < d.length; i++) {
                                    content += `<p>${d[i]}</p>`
                                }
                                that.$Modal.success({
                                    title: "消息",
                                    content
                                })
                            }
                            that.interface()
                        })
                        .catch(() => {
                            that.interface()
                        })
                    that.updateStopservice = false
                } else {
                    that.updateStopservice = true
                }
            })
        },
        updateStopserviceCancel() {
            this.$refs.updateStopserviceRef.resetFields()
            this.updateStopservice = false
        },
        configureCancel() {
            this.$refs.formConfigureValidate.resetFields()
            this.configureModal = false
        },
        sethostGroupList(page) {
            this.newHostGroup = this.hostGroupList.slice(page * 10, (page + 1) * 10)
            // console.log('this.newHostGroup',this.newHostGroup)
        },
        setClustersList(page) {
            this.newClustersList = this.hostgroupClustersList.slice(page * 10, (page + 1) * 10)
        },
        setPullList(page) {
            this.newPullList = this.pullList.slice(page * 10, (page + 1) * 10)
        },
        uploading() {
            const inputDOM = this.$refs.inputer;
[this.initValidate.file] = inputDOM.files
        },
        pushTimeChange(time) {
            this.configureValidate.pushTime = time
        },
        getBasic() {
            const that = this
            API.Basic({ clusterName: that.trim(that.clusterName), snapshotName: that.trim(that.snapshotName) }).then(
                res => {
                    that.dataList = res
                }
            )
        },
        stopService() {
            const that = this
            /* eslint-disable */
            API.stopService({
                clusterName: that.trim(that.clusterName),
                snapshotName: that.trim(that.snapshotName)
            }).then(res => {
                that.stopServiceList = res
            })
        },
        switchConfig() {
            let that = this
            /* eslint-disable */
            API.switchConfig({
                clusterName: that.trim(that.clusterName),
                snapshotName: that.trim(that.snapshotName)
            }).then(res => {
                that.switchConfigList = res
            })
        },
        hostgroups() {
            let that = this
            /* eslint-disable */
            API.hostgroups({
                clusterName: that.trim(that.clusterName),
                snapshotName: that.trim(that.snapshotName)
            }).then(res => {
                that.hostGroupList = res
                that.sethostGroupList(0)
            })
        },
        hostgroupClusters() {
            let that = this
            /* eslint-disable */
            API.hostgroupClusters({
                clusterName: that.trim(that.clusterName),
                snapshotName: that.trim(that.snapshotName)
            }).then(res => {
                that.hostgroupClustersList = res
                that.setClustersList(0)
            })
        },
        full() {
            let that = this
            API.full({ clusterName: that.trim(that.clusterName), snapshotName: that.trim(that.snapshotName) }).then(
                res => {
                    that.pullList = res
                    that.setPullList(0)
                }
            )
        },
        tabClick(name) {
            if (this.clusterName != "") {
                this.interface()
            }
        },
        clusterNameChage(name) {
            let that = this
            that.snapshotName = ""
            that.interface()
            that.getSnapshotList(name)
        },
        getSnapshotList(name) {
            if (name != "" && typeof name != "undefined") {
                this.$store.dispatch("zerodb/getSnapshotList", name)
            }
        },
        snapshotNameChage() {
            this.interface()
        },
        interface() {
            if (this.newTabName == "name1") {
                this.getBasic()
            }
            if (this.newTabName == "name2") {
                this.stopService()
            }
            if (this.newTabName == "name3") {
                this.switchConfig()
            }
            if (this.newTabName == "name4") {
                this.hostgroups()
            }
            if (this.newTabName == "name5") {
                this.hostgroups()
                this.hostgroupClusters()
            }
            if (this.newTabName == "name6") {
                this.hostgroupClusters()
                this.full()
            }
        },
        creatSnapshotOk() {
            let that = this
            that.$refs["formSnapshotValidate"].validate(valid => {
                if (valid) {
                    API.creatSnapshot({
                        clusterName: that.snapshotValidate.addClusterName,
                        snapshotName: that.trim(that.snapshotValidate.addSnapshotName)
                    }).then(res => {
                        if (res == null) {
                            that.$Modal.success({
                                title: "消息",
                                content: "创建成功!"
                            })
                        } else {
                            let d = res
                            let content = ""
                            for (let i = 0; i < d.length; i++) {
                                content = content + `<p>${d[i]}</p>`
                            }
                            that.$Modal.success({
                                title: "消息",
                                content: content
                            })
                        }
                        if (that.clusterName == that.snapshotValidate.addClusterName) {
                            that.$store.dispatch("zerodb/getSnapshotList", that.snapshotValidate.addClusterName)
                        }
                        that.$refs["formSnapshotValidate"].resetFields()
                    })
                    that.snapshotModal = false
                } else {
                    that.snapshotModal = true
                }
            })
        },
        creatSnapshotCancel() {
            this.$refs["formSnapshotValidate"].resetFields()
            this.snapshotModal = false
        },
        configureOk() {
            let that = this
            that.$refs["formConfigureValidate"].validate(valid => {
                if (valid) {
                    API.pushConfig({
                        clusterName: that.trim(that.configureValidate.configureClusterName),
                        snapshotName: that.trim(that.configureValidate.configureSnapshotName),
                        doTime: that.trim(that.configureValidate.pushTime)
                    }).then(res => {
                        if (res == null) {
                            that.$Modal.success({
                                title: "消息",
                                content: "推送配置成功!"
                            })
                        } else {
                            let d = res
                            let content = ""
                            for (let i = 0; i < d.length; i++) {
                                content = content + `<p>${d[i]}</p>`
                            }
                            that.$Modal.success({
                                title: "消息",
                                content: content
                            })
                        }
                    })
                    that.$refs["formConfigureValidate"].resetFields()
                    that.configureModal = false
                } else {
                    that.configureModal = true
                }
            })
        },
        initOk() {
            let that = this
            that.$refs["formInitValidate"].validate(valid => {
                if (valid) {
                    let formdata = new FormData()
                    formdata.append("file", that.initValidate.file)
                    formdata.append("clusterName", that.initValidate.initClusterName)
                    formdata.append("force", "0")
                    API.initConfig(formdata).then(res => {
                        if (res == null) {
                            that.$Modal.success({
                                title: "消息",
                                content: "初始化配置成功!"
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
                    })
                    that.$refs["formInitValidate"].resetFields()
                    that.initModal = false
                } else {
                    that.initModal = true
                }
            })
        },
        initCancel() {
            this.$refs["formInitValidate"].resetFields()
            this.initModal = false
        },
        exportOk() {
            let that = this
            that.$refs["formExportValidate"].validate(valid => {
                if (valid) {
                    API.exportConfig({
                        clusterName: that.exportValidate.exportClusterName,
                        snapshotName: that.exportValidate.exportSnapshotName,
                        exportName: that.trim(that.exportValidate.exportName)
                    }).then(res => {
                        window.location.href = res
                        that.$refs["formExportValidate"].resetFields()
                        that.exportModal = false
                    })
                } else {
                    that.exportModal = true
                }
            })
        },
        exportCancel() {
            this.$refs["formExportValidate"].resetFields()
            this.exportModal = false
        }
    },
    components: {
        Crumb,
        PaddingView,
        dataBase,
        subTable,
        hostGroupCluster
    },
    computed: mapGetters({
        snapshotList: "zerodb/snapshotList",
        clusterList: "zerodb/clusterList"
    }),
    created() {
        this.$store.dispatch("zerodb/getClusterList")
    }
}
</script>

<style scoped lang="less">
.configure-select {
    width: 100px;
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
.text-r {
    text-align: right;
}
.margin-t20 {
    margin-top: 20px;
}
.margin-b10 {
    margin-bottom: 10px;
}
.configure-w {
    width: 150px;
}
</style>
