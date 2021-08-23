<template>
    <div class="app-container">
        <el-card class="box-card">
  <div slot="header" class="clearfix">
    <span>任务集群</span>
    <el-button style="float: right; padding: 3px 0" @click.native="addFormVisible=true" type="text">新建任务集群</el-button>
  </div>
  <el-table
    :data="groups"
    style="width: 100%">
    <el-table-column
      prop="name"
      label="任务集群"
      width="200">
      <template slot-scope="scope">
         <el-tag  type="danger" size="medium" style="panding-right:3px;">{{scope.row.name}}</el-tag> 
      </template>
    </el-table-column>
    <el-table-column
      prop="remark"
      width="200"
      label="任务集群描述">
    </el-table-column>
    <el-table-column
      prop="clients"
      label="Clients">
      <template slot-scope="scope">
         <el-tag  style="margin-right:3px;" type="success" size="medium" v-for="ip in scope.row.clients" :key="ip">{{ip}}</el-tag>
      </template>
    </el-table-column>
  </el-table>
</el-card>
 <!--新增界面-->
    <el-dialog title="新增任务集群" :visible.sync="addFormVisible" :size="small" :close-on-click-modal="false">
      <el-form :model="addForm" label-width="100px" :rules="addFormRules" ref="addForm">
        <el-form-item label="集群名称" prop="name">
          <el-input v-model="addForm.name" auto-complete="off" maxlength="15" placeholder="请输入任务集群名称"></el-input>
        </el-form-item>
       
        <el-form-item label="备注" prop="remark">
             <el-input v-model="addForm.remark" auto-complete="off" placeholder="请输入任务集群备注信息"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false" size="mini">取消</el-button>
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading" size="mini">提交</el-button>
      </div>
    </el-dialog>


    </div>
</template>

<script>

import {findGroupDetailsList,addGroup}from '@/api/group'
export default {
    
    data(){
        return {
            groups:[],
            addFormVisible: false, // 新增界面是否显示
            addLoading: false,
            loading:false,
            addFormRules: {
              name: [
                { required: true, message: '请输入任务集群名称', trigger: 'blur' }
              ],
              remark: [
                { required: true, message: '请输入任务集群备注信息', trigger: 'blur' }
              ]
            },
            addForm: {
              name: '',
              remark:''
            }
          
        }
    },
    created(){
      this.searchGroupList()
    },
    methods:{

    /**
     * 添加任务集群
     */
    addSubmit(){
     this.$refs.addForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认要新增此任务集群信息吗？', '友情提示', {}).then(() => {
            this.addLoading = true
            let para = Object.assign({}, this.addForm)
            addGroup(para).then((res) => {
              this.addLoading = false
              this.$message({
                    message: "添加成功",
                    type: 'success'
                  })
              // reset the form 
              this.$refs['addForm'].resetFields()
              this.addFormVisible = false
              this.searchGroupList()
            }).catch((error)=>{
                 this.addLoading = false
            })
          })
        }
      })
   },
        /**
         * 查询任务集群列表
         */
        searchGroupList(){
          findGroupDetailsList({}).then(resp =>{
            this.groups = resp.data
        })
        }
    }
}
</script>