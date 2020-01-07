<template>
    <!-- eslint-disable -->
    <div>
        <Crumb></Crumb>
        <div class="container">
            <padding-view>
                <div class="search-bar clearfix margin-b20">
                    <div class="search-value fl">
                        <Select v-model="clusterName" @on-change="clusterNameChage" placeholder="群集名称">
                            <Option v-for="(item, k) in clusterList" :value="item" :key="k">{{ item }}</Option>
                        </Select>
                    </div>
                    <Button class="margin-l10" @click="modal1 = true">切换数据源</Button>
                    <Button class="margin-l10" @click="refresh">刷新表格</Button>
                </div>
                <f-table :data="newcolonyList">
                    <f-table-column title="ip" prop="ip"></f-table-column>
                    <f-table-column title="角色">
                        <template slot-scope="scope">
                            <div>Proxy</div>
                        </template>
                    </f-table-column>
                    <f-table-column title="CPU" prop="cpuLoad"></f-table-column>
                    <f-table-column title="Memory" prop="memLoad"></f-table-column>
                    <f-table-column title="Loadavg" prop="loadAvg"></f-table-column>
                    <f-table-column title="状态" prop="status"></f-table-column>
                    <f-table-column title="配置版本" prop="ConfVersion"></f-table-column>
                    <f-table-column title="操作">
                        <template slot-scope="scope">
                            <inline-button @on-tap="detail(scope)">下线</inline-button>
                        </template>
                    </f-table-column>
                </f-table>
                <Page
                    class="page-bar"
                    :total="colonyList.length"
                    :page-size="10"
                    @on-change="pageOnChange"
                    show-total
                ></Page>
            </padding-view>
        </div>
        <Modal v-model="unregisterModel" title="下线" @on-ok="unregisterOk">
            <Row :gutter="16" class="margin-b20">
                <Col span="4" class="text-r">集群名称：</Col>
                <Col span="18">
                    <Select v-model="clusterName" placeholder="群集名称" disabled>
                        <Option v-for="(item, k) in clusterList" :value="item" :key="k">{{ item }}</Option>
                    </Select>
                </Col>
            </Row>
            <Row :gutter="16" class="margin-b20">
                <Col span="4" class="text-r">代理服务器：</Col>
                <Col span="18">
                    <Input v-model="unregister.address" disabled />
                </Col>
            </Row>
            <Row :gutter="16" class="margin-b20">
                <Col span="4" class="text-r">原因：</Col>
                <Col span="18">
                    <Input v-model="unregister.reason" placeholder="原因" />
                </Col>
            </Row>
        </Modal>

        <Modal v-model="modal1" title="切换数据源">
            <Form ref="formColonyValidate" :model="colony" label-position="left" :label-width="190">
                <FormItem
                    label="群集名称"
                    prop="colonyClusterName"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select
                        v-model="colony.colonyClusterName"
                        @on-change="colonyClusterNameChage"
                        placeholder="群集名称"
                    >
                        <Option v-for="(item, k) in clusterList" :value="item" :key="k">{{ item }}</Option>
                    </Select>
                </FormItem>
                <FormItem
                    label="数据库组"
                    prop="colonyhostGroup"
                    :rules="{ required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select v-model="colony.colonyhostGroup" placeholder="数据库组">
                        <Option v-for="(item, k) in hostGroupList" :value="item.name" :key="k">{{ item.name }}</Option>
                    </Select>
                </FormItem>
                <FormItem
                    label="切出序号"
                    prop="from"
                    :rules="{ type: 'number', required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select v-model="colony.from" placeholder="切出序号">
                        <Option v-for="(item, k) in setWrite(colony.colonyhostGroup)" :value="k" :key="k">{{
                            k
                        }}</Option>
                    </Select>
                </FormItem>
                <FormItem
                    label="切入序号"
                    prop="to"
                    :rules="{ type: 'number', required: true, message: '不能为空', trigger: 'blur' }"
                >
                    <Select v-model="colony.to" placeholder="切入序号">
                        <Option v-for="(item, k) in setWrite(colony.colonyhostGroup)" :value="k" :key="k">{{
                            k
                        }}</Option>
                    </Select>
                </FormItem>
                <FormItem label="原因" prop="reason" :rules="{ required: true, message: '不能为空', trigger: 'blur' }">
                    <Input v-model="colony.reason" placeholder="原因" />
                </FormItem>
            </Form>
            <div slot="footer">
                <Button type="text" size="large" @click="toggleCancel">取消</Button>
                <Button type="primary" size="large" @click="toggleOk">确定</Button>
            </div>
        </Modal>
    </div>
</template>

<script>
import Crumb from "@components/crumb"
import PaddingView from "@components/padding-view"
import fTable from "@components/fTable/fTable"
import fTableColumn from "@components/fTable/fTableColumn"
import InlineButton from "@components/button/inline-button"
import { mapGetters } from "vuex"
import API from "../../../__BASE__/api"

export default {
    data() {
        return {
            dataList: {},
            clusterName: "",
            hostGroup: "",
            modal1: false,
            hostGroupList: [],
            colonyList: [],
            unregisterModel: false,
            unregister: {
                address: "",
                reason: ""
            },
            newcolonyList: [],
            colony: {
                colonyClusterName: "",
                colonyhostGroup: "",
                from: "",
                to: "",
                reason: ""
            }
        }
    },
    components: {
        Crumb,
        PaddingView,
        fTable,
        fTableColumn,
        InlineButton
    },
    methods: {
        refresh() {
            this.setColonyList()
        },
        detail(item) {
            this.unregisterModel = true
            this.unregister.address = item.ip
        },
        unregisterOk() {
            const that = this
            that.unregister.clusterName = that.clusterName
            API.unregister(that.unregister).then(res => {
                if (res == null) {
                    that.$Modal.success({
                        title: "消息",
                        content: "下线成功!"
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
                that.setColonyList()
            })
            that.$store.dispatch("zerodb/getClusterList")
        },
        setWrite(data) {
            if (data === "" || typeof data === "undefined") {
                return []
            }
            let w
            const hgl = this.hostGroupList
            for (let i = 0; i < hgl.length; i++) {
                if (hgl[i].name === data) {
                    w = hgl[i].write
                }
            }
            return w.split(",")
        },
        pageOnChange(page) {
            const p = page - 1
            this.newcolonyList = this.colonyList.slice(p * 10, (p + 1) * 10)
        },
        pageSizeChange() {},
        setColonyList() {
            const that = this
            API.proxyStatus({ clusterName: that.clusterName })
                .then(res => {
                    that.colonyList = res
                    that.pageOnChange(1)
                })
                .catch(() => {
                    that.colonyList = []
                    that.pageOnChange(1)
                })
        },
        clusterNameChage() {
            this.setColonyList()
        },
        colonyClusterNameChage() {
            const that = this
            if (that.colony.colonyClusterName) {
                API.hostgroups({ clusterName: that.colony.colonyClusterName }).then(res => {
                    that.hostGroupList = res
                })
            }
        },
        toggleOk() {
            const that = this
            this.$refs.formColonyValidate.validate(valid => {
                if (valid) {
                    /* eslint-disable */
                    API.switchProxy({
                        clusterName: that.colony.colonyClusterName,
                        hostGroup: that.colony.colonyhostGroup,
                        from: that.colony.from,
                        to: that.colony.to,
                        reason: that.colony.reason
                    }).then(res => {
                        if (res == null) {
                            this.$Modal.success({
                                title: "消息",
                                content: "切换成功!"
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
                        that.setColonyList()
                    })
                    that.modal1 = false
                    that.$refs["formColonyValidate"].resetFields()
                } else {
                    that.modal1 = true
                }
            })
        },
        toggleCancel() {
            this.$refs["formColonyValidate"].resetFields()
            this.modal1 = false
        }
    },
    computed: mapGetters({
        clusterList: "zerodb/clusterList"
    }),
    mounted() {
        this.$store.dispatch("zerodb/getClusterList")
    }
}
</script>

<style scoped>
.container {
    padding: 30px;
}
.margin-l10 {
    margin-left: 10px;
}
.search-value {
    width: 150px;
}
</style>
