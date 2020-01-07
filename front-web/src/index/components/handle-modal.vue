<!--
    流程发起modal  zipai at 2018/09/09
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
                {{ handleParams.title }}
            </p>
            <div class="content">
                <Form
                    ref="formValidate"
                    :model="formValidate"
                    :rules="ruleValidate"
                    :label-width="100"
                    class="margin_r40"
                >
                    <div v-for="(item, index) in baseData.fields" :key="index">
                        <!-- 单行文本框 text-->
                        <FormItem v-if="item.type === 'text'" :label="item.name + '：'" :prop="item.id">
                            <Input
                                class="width300"
                                :disabled="item.read_only"
                                v-model="formValidate[item.id]"
                                :autocomplete="'off'"
                                :placeholder="item.placeholder"
                            ></Input>
                        </FormItem>
                        <!-- 多行文本框 multi-line-text-->
                        <FormItem v-if="item.type === 'multi-line-text'" :label="item.name + '：'" :prop="item.id">
                            <Input
                                class="width300"
                                type="textarea"
                                :disabled="item.read_only"
                                v-model="formValidate[item.id]"
                                :autocomplete="'off'"
                                :placeholder="item.placeholder"
                                :maxlength="getMaxLen(item.params)"
                            ></Input>
                        </FormItem>
                        <!-- 数字 integer-->
                        <FormItem v-if="item.type === 'integer'" :label="item.name + '：'" :prop="item.id">
                            <InputNumber
                                class="width300"
                                :min="0"
                                :disabled="item.read_only"
                                v-model="formValidate[item.id]"
                                :placeholder="item.placeholder"
                                :maxlength="getMaxLen(item.params)"
                            ></InputNumber>
                        </FormItem>
                        <!-- 线标题 headline-->
                        <FormItem v-if="item.type === 'headline'" class="headline" :label="item.name + '：'"></FormItem>
                        <!--多选 Checkbox -->
                        <FormItem v-if="item.type === 'boolean'" class="checkbox">
                            <Checkbox v-model="formValidate[item.id]">{{ item.name }}</Checkbox>
                        </FormItem>
                        <!-- 下拉选择  dropdown-->
                        <FormItem v-if="item.type === 'dropdown'" :label="item.name + '：'" :prop="item.id">
                            <Select
                                class="width300"
                                :disabled="item.read_only"
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
                        <FormItem v-if="item.type === 'date'" :label="item.name + '：'" :prop="item.id">
                            <DatePicker
                                class="width300"
                                type="date"
                                :disabled="item.read_only"
                                :value="formValidate[item.id]"
                                @on-change="dataChange($event, item.id)"
                                placeholder="请输入时间"
                            ></DatePicker>
                        </FormItem>
                        <!-- people 一期先按下拉搜索选择做-->
                        <FormItem v-if="item.type === 'people'" :label="item.name + '：'" :prop="item.id">
                            <Select
                                class="width300"
                                :disabled="item.read_only"
                                :readonly="item.read_only"
                                :placeholder="item.placeholder"
                                v-model="formValidate[item.id]"
                                filterable
                                remote
                                @on-query-change="remoteMethod($event, item)"
                                :loading="remoteLoading"
                                @on-change="setUserVal($event, item.id)"
                            >
                                <Option v-for="item1 in item.users" :key="item1.id" :value="item1.id">{{
                                    item1.fullName
                                }}</Option>
                            </Select>
                        </FormItem>
                    </div>
                </Form>

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
            <div slot="footer" class="footer">
                <div v-if="baseData.outcomes != ''">
                    <Button
                        type="primary"
                        class="modal-btn"
                        v-for="item in baseData.outcomes"
                        :key="item.id"
                        @click="handleSubmit('formValidate', item.name)"
                        >{{ item.name }}</Button
                    >
                </div>
                <div v-else>
                    <Button type="primary" class="modal-btn" @click="handleSubmit('formValidate', '完成')">完成</Button>
                </div>
                <Button
                    v-if="!handleParams.hasform"
                    type="primary"
                    class="modal-btn"
                    @click="handleSubmit('formValidate', '完成')"
                    >完成</Button
                >
            </div>
        </Modal>
    </div>
</template>
<script>
import InlineTime from "@components/inline-time"
import InlineLabel from "@components/inline-label"
import API from "@api"
export default {
    props: ["handleParams"],
    data() {
        return {
            baseData: {},
            formValidate: {},
            ruleValidate: {},
            Operations: [],
            params: {
                form_id: null,
                outcome: "",
                id: this.handleParams.id
            },
            userArr: [],
            remoteLoading: false
        }
    },
    components: {
        InlineLabel,
        InlineTime
    },
    methods: {
        cancel() {
            this.$emit("on-cancel")
        },
        emitOk() {
            let fromVal = JSON.parse(JSON.stringify(this.formValidate))
            this.userArr.forEach(v => {
                let key = v.key
                let val = v.value
                fromVal[key] = val
            })
            this.params.values = fromVal
            this.$emit("on-ok", this.params)
        },
        dataChange(time, key) {
            this.formValidate[key] = time
        },
        getMaxLen(v) {
            if (v && v.maxLength) {
                return Number(v.maxLength)
            }
            return 10000
        },
        handleSubmit(name, btnName) {
            this.params.outcome = btnName
            if (!this.handleParams.hasform) {
                //不存在表单的情况
                this.$emit("on-ok", this.params)
            } else {
                this.$refs[name].validate(valid => {
                    if (valid) {
                        this.emitOk()
                    } else {
                        this.$Message.error("信息不完整!")
                    }
                })
            }
        },
        uniq(key) {
            this.userArr.forEach((v, i) => {
                if (v.key === key) {
                    v.value = params
                    return true
                }
            })
        },
        toArr(datas, key) {
            let params = {
                key: key,
                value: {
                    id: datas.id,
                    firstName: datas.firstName,
                    lastName: datas.lastName,
                    email: datas.email
                }
            }
            if (!this.uniq(key)) {
                this.userArr.push(params)
            }
        },
        getuser(v, val) {
            let self = this
            API.getUsers({ filter: val })
                .then(data => {
                    v.users = data.data
                    if (data.data.length == 1) {
                        let datas = data.data[0]
                        self.toArr(datas, v.id)
                    }
                })
                .catch(e => {
                    console.log(e)
                })
        },
        setUserVal(v, key) {
            let self = this
            API.getUsers({ filter: v })
                .then(data => {
                    if (data && data.data && data.data[0]) {
                        let datas = data.data[0]
                        self.toArr(datas, key)
                    }
                })
                .catch(e => {})
        },
        remoteMethod(query, item) {
            this.remoteLoading = true
            API.getUsers({ filter: query })
                .then(data => {
                    this.remoteLoading = false
                    item.users = data.data
                })
                .catch(e => {
                    this.remoteLoading = false
                })
        },
        /**
         * 获取流程流水
         */
        flowOperations() {
            let self = this
            API.flowOperations({ id: self.handleParams.p_i_id })
                .then(data => {
                    self.Operations = data
                })
                .catch(err => {
                    console.log(err)
                })
        },
        getData() {
            // 查询表单数据
            let self = this
            API.getTaskForm({ id: self.handleParams.id })
                .then(data => {
                    self.baseData = data
                    self.params.form_id = data.id
                    data.fields.forEach((v, i) => {
                        self.formValidate[v.id] = v.value
                        self.ruleValidate[v.id] = []
                        if (v.type === "people") {
                            self.$set(v, "users", [])
                            self.getuser(v, v.value)
                        }
                        if (v.type === "integer") {
                            self.ruleValidate[v.id].push({
                                required: v.required,
                                message: v.required ? "此项为必填项" : "",
                                trigger: "change",
                                type: "number"
                            })
                        } else {
                            self.ruleValidate[v.id].push({
                                required: v.required,
                                message: v.required ? "此项为必填项" : "",
                                trigger: "change"
                            })
                        }
                        if (v.type === "multi-line-text" || v.type === "text") {
                            let regRule = data.fields[i].params
                            if (regRule && regRule.regexPattern) {
                                let reg = new RegExp(`^${regRule.regexPattern}$`)
                                console.log(reg, "reg")
                                self.ruleValidate[v.id].push({
                                    trigger: "change",
                                    validator(rule, val, callback, source, options) {
                                        if (val && regRule.minLength && val.length < regRule.minLength) {
                                            callback("输入错误！输入长度至少为" + regRule.minLength + "位")
                                        } else if (val && !reg.test(val)) {
                                            callback("输入错误！输入规则不匹配")
                                        } else {
                                            callback()
                                        }
                                    }
                                })
                            }
                        }
                    })
                })
                .catch(e => {
                    console.log(e)
                })
        }
    },
    created() {
        if (this.handleParams.hasform) {
            this.getData()
        }
        this.flowOperations()
    }
}
</script>
<style lang="less" scoped>
.modal /deep/ .ivu-modal-footer {
    border-top: 0;
}
.modal /deep/ .headline.ivu-form-item {
    margin-bottom: 0px;
}
.modal /deep/ .checkbox.ivu-form-item {
    margin-bottom: 10px;
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
.footer {
    text-align: center;
    padding-bottom: 15px;
    .modal-btn {
        margin: 10px 20px;
    }
}
</style>
