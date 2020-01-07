<!--
    员工详情 模态框  zipai at 2018/6/5
-->
<template>
    <div class="employees-detail">
        <Modal class="modal employees-detail-modal" width="500px" :mask-closable="false" :value="true" @on-cancel="cancel">
            <h3 slot="header" class="title">{{title}}</h3>
            <div class="content">
							<div  v-for="(firstItem,index) in info" :key='index'>
								<div v-for="baseItem in firstItem" :key='baseItem.id'>
                    <InfoTitle :title="baseItem.category"></InfoTitle>
                    <div v-for="item in baseItem.fieldList" :key='item.id'>
                        <BaseInput :title="item.name+'：'"  v-if="item.name ==='手机号码'">
                            <p slot="input" class="info">
                                <span class="phone" v-if="isShow">{{item.value}}</span>
                                <span class="phone" v-else>{{item.value | formatPhone}}</span>
                                <!-- <Button @click="switchover" size="small" class="btn" v-text="text"></Button> -->
                            </p>
                        </BaseInput>
                        <BaseInput v-else :title="item.name+'：'">
                            <span slot="input" class="info">{{item.value}}</span>
                        </BaseInput>
                    </div>
                </div>
							</div>
            </div>
            <div slot="footer" class="footer"></div>
        </Modal>
    </div>
</template>
<script>
import API from '@api'
import BaseInput from "./base-input";
import InfoTitle from "@components/info-title";
export default {
    props: ["title", "id",'enabled'],
    data() {
        return {
            value: null,
            isShow:false,
            text:'显示',
            params: {
                userId: this.id,
                enabled:this.enabled==='enabled'?false:true,
            },
            info: [
                {
                    category: "个人信息",
                    fieldList: [
                        {
                            id: "8516fd1f011040959cb6654034b86359",
                            name: "花名",
                            displayName: null,
                            value: "XXXXXXX",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "f69511101f74415b900164adb615bb62",
                            name: "姓名",
                            displayName: null,
                            value: "XXXXXXX",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "b929e68bc3494fddad5ad9cb835f29e2",
                            name: "性别",
                            displayName: null,
                            value: "",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "539a27b4dcd94b52b76914b2f81b2f53",
                            name: "手机号码",
                            displayName: null,
                            value: "133****3421",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        }
                    ]
                },
                {
                    category: "公司信息",
                    fieldList: [
                        {
                            id: "ed7a0a6b889743d39a9a6eb818f45fb5",
                            name: "工号",
                            displayName: null,
                            value: "001",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "afd1966a1b2142e983a697099730ab92",
                            name: "部门",
                            displayName: null,
                            value: "总裁部",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "cafe41a5b0264c4dbcae1be8fc3ddb68",
                            name: "职位",
                            displayName: null,
                            value: "",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "2971b5c6d7624e32bd0c8dd033854040",
                            name: "级别",
                            displayName: null,
                            value: "",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "b2872da434764e8c925697ba82ff9d2c",
                            name: "电子邮箱",
                            displayName: null,
                            value: "zongcai@qq.com",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        },
                        {
                            id: "f03f27b4b6254fea925464cd525b15a7",
                            name: "员工类型",
                            displayName: null,
                            value: "",
                            categoryId: null,
                            sort: 0,
                            typeId: null,
                            type: null,
                            extra: null,
                            notNull: false,
                            isOnly: null,
                            visible: false,
                            modify: false
                        }
                    ]
                }
            ]
        };
    },
    filters:{
        formatPhone(v){
            return v.substr(0,3)+"****"+v.substr(7);
        }
    },
    components: { BaseInput, InfoTitle },
    methods: {
        init() {
            this.getUserInfo();
        },
        getUserInfo() {
            API.getApplyUserInfo({id:this.id}).then(data => {
							console.log(data);
                this.info = data;
            }).catch((e)=>{
							console.log(e);
						});
        },
        cancel() {
            this.$emit("on-cancel");
        },
        switchover(){
            this.isShow = !this.isShow;
            this.text = this.isShow?'隐藏':'显示'
        }
    },
    created() {
        console.log(this.id)
        this.init();
    },
    mounted() {}
};
</script>
<style lang="less" scoped>
.modal /deep/ .ivu-modal-body{
	min-height: 360px;
  max-height: 700px;
  overflow-y: auto;
}
.input_wrapper{
	padding-left: 60px;
}
.title {
    font-size: 14px;
    color: #333333;
    font-weight: bold;
    text-align: center;
}
.content {
    // width: 434px;
    // margin: 0 auto;
    // padding: 0 0 20px 25px;
    .info {
        display: block;
        height: 30px;
        line-height: 30px;
        .phone{
            display: inline-block;
            width: 100px;
        }
        .btn{
            display: inline-block;
        }
    }
}
.footer {
    text-align: center;
    .btn {
        margin: 0 40px;
    }
}
</style>


