<!--
    详情展示modal  zipai at 2018/09/09
	@props:
     header: header显示文字内容，
	 id:根据id查取数据
-->
<template>
    <div class="details-modal">
        <Modal
            :value="true"
            @on-cancel="cancel"
            :mask-closable="false"
            class="modal"
            class-name="vertical-center-modal"
            :width="500"
        >
            <p class="header" slot="header">
                {{ title }}
            </p>
            <div class="content">
                <!-- <inline-desc title='花名：'>小小酥</inline-desc> -->
                <div class="form" v-if="fromData">
                    <Form ref="formValidate" :model="formValidate" :label-width="100" class="margin_r40">
                        <div v-for="(item, index) in fromData.fields" :key="index">
                            <!-- 单行文本框 text-->
                            <FormItem v-if="item.type === 'text'" :label="item.name + '：'">
                                <Input
                                    class="width300"
                                    disabled
                                    v-model="formValidate[item.id]"
                                    :autocomplete="'off'"
                                    :placeholder="item.placeholder"
                                ></Input>
                            </FormItem>
                            <!-- 多行文本框 multi-line-text-->
                            <FormItem v-if="item.type === 'multi-line-text'" :label="item.name + '：'">
                                <Input
                                    class="width300"
                                    type="textarea"
                                    disabled
                                    v-model="formValidate[item.id]"
                                    :autocomplete="'off'"
                                    :placeholder="item.placeholder"
                                    :maxlength="getMaxLen(item.params)"
                                ></Input>
                            </FormItem>
                            <!-- 数字 integer-->
                            <FormItem v-if="item.type === 'integer'" :label="item.name + '：'">
                                <InputNumber
                                    class="width300"
                                    :min="0"
                                    disabled
                                    v-model="formValidate[item.id]"
                                    :placeholder="item.placeholder"
                                    :maxlength="getMaxLen(item.params)"
                                ></InputNumber>
                            </FormItem>
                            <!-- 线标题 headline-->
                            <FormItem
                                v-if="item.type === 'headline'"
                                class="headline"
                                :label="item.name + '：'"
                            ></FormItem>
                            <!--多选 Checkbox -->
                            <FormItem v-if="item.type === 'boolean'" class="checkbox">
                                <Checkbox v-model="formValidate[item.id]" disabled>{{ item.name }}</Checkbox>
                            </FormItem>

                            <!-- 下拉选择  dropdown-->
                            <FormItem v-if="item.type === 'dropdown'" :label="item.name + '：'">
                                <Select
                                    class="width300"
                                    disabled
                                    :placeholder="item.placeholder"
                                    v-model="formValidate[item.id]"
                                >
                                    <Option
                                        v-for="(selectItem, index) in item.options"
                                        :key="index"
                                        :value="selectItem.name"
                                        >{{ selectItem.name }}</Option
                                    >
                                </Select>
                            </FormItem>
                            <!-- 日期 date-->
                            <FormItem v-if="item.type === 'date'" :label="item.name + '：'">
                                <DatePicker
                                    class="width300"
                                    type="date"
                                    readonly
                                    disabled
                                    :value="formValidate[item.id]"
                                    @on-change="dataChange($event, item.id)"
                                    placeholder="请输入时间"
                                ></DatePicker>
                            </FormItem>
                            <!-- people 一期先按下拉搜索选择做-->
                            <FormItem v-if="item.type === 'people'" :label="item.name + '：'">
                                <Select
                                    class="width300"
                                    disabled
                                    :placeholder="item.placeholder"
                                    v-model="formValidate[item.id]"
                                >
                                    <Option v-for="item1 in item.users" :key="item1.id" :value="item1.id">{{
                                        item1.fullName
                                    }}</Option>
                                </Select>
                            </FormItem>
                        </div>
                    </Form>
                </div>
                <p class="line"></p>
                <Timeline class="timeline">
                    <TimelineItem v-for="(item, index) in Operations" :key="index">
                        <inline-time :title="item.name">
                            <p v-for="(item2, index) in item.user_operations" :key="index">
                                {{ item2.start_time + "　" + item2.operation }}
                            </p>
                        </inline-time>
                    </TimelineItem>
                </Timeline>
            </div>
            <div slot="footer"></div>
        </Modal>
    </div>
</template>
<script>
import API from "@api"
import InlineDesc from "@components/inline-desc"
import InlineTime from "@components/inline-time"
export default {
    props: ["id", "title", "hasform"],
    data() {
        return {
            formValidate: {},
            fromData: null,
            Operations: [],
            users: []
        }
    },
    components: {
        InlineDesc,
        InlineTime
    },
    methods: {
        init() {
            if (this.hasform) {
                this.getFlowDetail(this.id)
            }
            this.flowOperations(this.id)
        },
        formInit(data) {
            let self = this
            data.forEach((v, i) => {
                if (v.type === "people") {
                    self.getuser(v, v.value)
                }
                if (v.type === "integer") {
                    v.value = Number(v.value)
                }
                self.formValidate[v.id] = v.value
            })
        },
        getMaxLen(v) {
            if (v && v.maxLength) {
                return Number(v.maxLength)
            }
            return 10000
        },
        cancel() {
            this.$emit("on-cancel")
        },
        /**
         * 获取流程实例详情
         */
        getFlowDetail(id) {
            let self = this
            API.getProcessDetail({ id: id }).then(data => {
                self.fromData = data.form
                self.formInit(data.form.fields)
            })
        },
        /**
         * 获取流程流水
         */
        flowOperations(id) {
            let self = this
            API.flowOperations({ id: id }).then(data => {
                self.Operations = data
            })
        },
        getuser(v, val) {
            API.getUsers({ filter: val })
                .then(data => {
                    this.$set(v, "users", data.data)
                })
                .catch(e => {})
        }
    },
    created() {
        this.init()
    }
}
</script>
<style lang="less" scoped>
.modal /deep/ .ivu-modal-footer {
    border-top: 0;
}
.timeline /deep/ .ivu-timeline-item-content {
    left: -10px;
}
.header {
    text-align: center;
}
.content {
    padding: 0 8px;
}
.width300 {
    width: 300px;
}
.line {
    margin: 20px 0;
    border-bottom: 1px solid #ececec;
}
.timeline {
    padding-left: 30px;
}
</style>
