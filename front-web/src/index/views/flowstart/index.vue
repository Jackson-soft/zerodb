<template>
    <div>
        <Crumb></Crumb>
        <div class="container">
            <div class="cardWrap">
                <min-card>
                    <span>待办事宜</span><br /><span>{{ tasks.task_todo_count }}</span>
                </min-card>
                <min-card>
                    <span>已办事宜</span><br /><span>{{ tasks.task_done_count }}</span>
                </min-card>
                <min-card>
                    <span>我发起的</span><br /><span>{{ tasks.process_instances_initiated_count }}</span>
                </min-card>
            </div>
            <padding-view>
                <div class="item" v-for="item in definitions" :key="item.process_definition_category_id">
                    <inside-title :text="item.process_definition_category_name"></inside-title>
                    <div class="item-cont">
                        <min-card class="card " v-for="secondItem in item.process_definitions" :key="secondItem.id">
                            <card-inside :text="secondItem.name" @on-tap="flowStart(secondItem)"></card-inside>
                        </min-card>
                    </div>
                </div>
            </padding-view>
        </div>
        <flowapply-modal
            v-if="isShowFlowApply"
            :id="flowId"
            :hasform="hasform"
            :header="flowName"
            @on-cancel="applyCancel"
            @on-ok="applyOk"
        ></flowapply-modal>
    </div>
</template>
<script>
import API from "@api"
import Crumb from "@components/crumb"
import PaddingView from "@components/padding-view"
import InsideTitle from "@components/inside-title"
import MinCard from "../../components/mini-card"
import CardInside from "../../components/card-inside"
import FlowapplyModal from "../../components/flowapply-modal"

export default {
    data() {
        return {
            isShowFlowApply: false,
            tasks: {
                process_instances_initiated_count: 200,
                task_done_count: 200,
                task_todo_count: 200
            },
            flowName: "",
            flowId: "",
            hasform: true,
            definitions: []
        }
    },
    components: {
        Crumb,
        PaddingView,
        MinCard,
        InsideTitle,
        CardInside,
        FlowapplyModal
    },
    methods: {
        /**
         * 事件绑定
         */
        init() {
            this.flowOverview()
            this.flowDefinitions()
        },
        flowStart(item) {
            this.isShowFlowApply = true
            this.flowId = item.id
            this.flowName = item.name
            this.hasform = item.start_form_defined
        },
        applyCancel() {
            this.isShowFlowApply = false
        },
        applyOk(v) {
            console.log(v, "成功申请")
            let self = this
            API.flowStart(v)
                .then(data => {
                    console.log(data)
                    self.isShowFlowApply = false
                    self.$Message.success("操作成功!")
                    self.flowOverview()
                })
                .catch(e => {})
        },
        /**
         * 数据请求
         */
        //流程概况
        flowOverview() {
            let self = this
            // API.flowOverview().then(data => {
            //   ;(self.tasks.process_instances_initiated_count = data.process_instances_initiated_count),
            //     (self.tasks.task_done_count = data.task_done_count),
            //     (self.tasks.task_todo_count = data.task_todo_count)
            // })
        },
        flowDefinitions() {
            let self = this
            // API.flowDefinitions().then(data => {
            //   self.definitions = data
            // })
        }
    },
    created() {
        this.init()
    }
}
</script>
<style scoped lang="less">
.container {
    padding: 30px;
}
.cardWrap {
    margin-bottom: 30px;
    span {
        font-weight: bold;
        font-size: 14px;
        color: #3c4144;
    }
}
.item {
    .item-cont {
        padding: 30px 0;
    }
}
</style>
