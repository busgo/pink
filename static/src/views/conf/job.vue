<template>
    <div class="app-container">

        <el-card class="box-card">
  <div slot="header" class="clearfix">
    <span>任务配置列表</span>
    <el-button style="float: right; padding: 3px 0" type="text" @click.native="addFormVisible=true">新建任务</el-button>
  </div>

<el-form :inline="true" :model="searchForm" class="demo-form-inline" size="small">
  <el-form-item label="任务集群">
    <el-select v-model="searchForm.group" style="width:100%"   clearable placeholder="请选择集群">
      <el-option v-for="item in groups" :key="item.name" :label="item.remark" :value="item.name"></el-option>
    </el-select>
  </el-form-item>
  <el-form-item>
    <el-button type="primary" @click="onSearchSubmit" icon="el-icon-search" round>查 询</el-button>
  </el-form-item>
</el-form>
  <el-table
    :data="jobs"
    style="width: 100%">
    <el-table-column
      fixed
      prop="id"
      label="ID"
      width="200">
    </el-table-column>
    <el-table-column
      prop="name"
      label="任务名称"
      width="150">
    </el-table-column>
    <el-table-column
      prop="group"
      label="集群"
      width="120">
    </el-table-column>
    <el-table-column
      prop="cron"
      label="Cron 表达式"
      width="150">
    </el-table-column>
    <el-table-column
      prop="state"
      label="状态"
      width="50">
      <template slot-scope="scope">
          <el-tag v-if="scope.row.state==1" type="success" size="mini">启 用</el-tag>
          <el-tag v-if="scope.row.state==2" type="danger" size="mini">关 闭</el-tag>
      </template>
    </el-table-column>
    <el-table-column
      prop="target"
      label="Target"
      width="350">
    </el-table-column>
     <el-table-column
      prop="param"
      label="参数"
      width="200">
    </el-table-column>
     <el-table-column
      prop="mobile"
      label="手机号码"
      width="150">
    </el-table-column>
     <el-table-column
      prop="create_time"
      label="创建时间"
      width="200">
    </el-table-column>
     <el-table-column
      prop="update_time"
      label="修改时间"
      width="200">
    </el-table-column>
    <el-table-column
      fixed="right"
      label="操作" width="350">
      <template slot-scope="scope">
        <el-button type="warning" size="mini" @click="handleExecute(scope.$index, scope.row)" round>立即执行</el-button>
        <el-button type="success" size="mini"  @click="handleEdit(scope.$index, scope.row)" round>编 辑</el-button>
         <el-button type="primary" size="mini" @click="handleView(scope.$index, scope.row)" round>详 情</el-button>
         <el-button type="danger" size="mini" @click="handleDelete(scope.$index, scope.row)" round>删 除</el-button>
      </template>
    </el-table-column>
  </el-table>
</el-card>


<!--新增界面-->
    <el-dialog title="新增任务" :visible.sync="addFormVisible" size="small" :close-on-click-modal="false">
      <el-form :model="addForm" label-width="100px" :rules="addFormRules" ref="addForm">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="addForm.name" auto-complete="off" maxlength="15"></el-input>
        </el-form-item>
        <el-form-item label="任务集群" prop="group">
      <el-select v-model="addForm.group" placeholder="请选择任务集群" style="width:100%;">
          <el-option
            v-for="item in groups"
            :key="item.name"
            :label="item.remark"
            :value="item.name">
          </el-option>
        </el-select>
        </el-form-item>
        <el-form-item label="Cron表达式" prop="cron">
             <el-input v-model="addForm.cron" auto-complete="off" placeholder="请输入Cron表达式"></el-input>
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-select v-model="addForm.state" placeholder="请选择任务状态" style="width:100%;">
          <el-option label="启用" :value="1"></el-option>
          <el-option label="关闭" :value="2"></el-option>
        </el-select>
        </el-form-item>
         <el-form-item label="Target" prop="target">
            <el-input v-model="addForm.target" placeholder="请输入Target" auto-complete="off"></el-input>
        </el-form-item>
         <el-form-item label="任务参数" prop="param">
            <el-input v-model="addForm.param" placeholder="请输入任务参数" auto-complete="off"></el-input>
        </el-form-item>
          <el-form-item label="手机号码" prop="mobile">
            <el-input v-model="addForm.mobile" placeholder="请输入手机号码" auto-complete="off"></el-input>
          </el-form-item>
            <el-form-item label="备注" prop="remark">
            <el-input v-model="addForm.remark" placeholder="请输入任务备注" auto-complete="off"></el-input>
          </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addFormVisible = false" size="mini">取消</el-button>
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading" size="mini">提交</el-button>
      </div>
    </el-dialog>



      <!--编辑界面-->
    <el-dialog title="修改任务配置" :visible.sync="editFormVisible" :close-on-click-modal="false">
      <el-form :model="editForm" label-width="100px" :rules="editFormRules" ref="editForm">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="editForm.name" auto-complete="off" maxlength="15"></el-input>
        </el-form-item>
        <el-form-item label="任务集群" prop="group">
      <el-select v-model="editForm.group" placeholder="请选择任务集群" style="width:100%;" disabled>
          <el-option
            v-for="item in groups"
            :key="item.name"
            :label="item.remark"
            :value="item.name">
          </el-option>
        </el-select>
        </el-form-item>
        <el-form-item label="Cron表达式" prop="cron">
             <el-input v-model="editForm.cron" auto-complete="off" placeholder="请输入Cron表达式" ></el-input>
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-select v-model="editForm.state" placeholder="请选择任务状态" style="width:100%;">
            <el-option label="开启" :value="1"></el-option>
            <el-option label="关闭" :value="2"></el-option>
        </el-select>
        </el-form-item>
         <el-form-item label="目标任务" prop="target">
            <el-input v-model="editForm.target" placeholder="请输入目标任务" auto-complete="off"></el-input>
        </el-form-item>
         <el-form-item label="任务参数" prop="param">
            <el-input v-model="editForm.param" placeholder="请输入任务参数" auto-complete="off"></el-input>
        </el-form-item>
          <el-form-item label="手机号码" prop="mobile">
            <el-input v-model="editForm.mobile" placeholder="请输入手机号码" auto-complete="off"></el-input>
          </el-form-item>
            <el-form-item label="备注" prop="remark">
            <el-input v-model="editForm.remark" placeholder="请输入任务备注" auto-complete="off"></el-input>
          </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editFormVisible = false" size="mini">取消</el-button>
        <el-button type="primary" @click.native="editSubmit" :loading="editLoading" size="mini">提交</el-button>
      </div>
    </el-dialog>

<!--查看界面-->
    <el-dialog title="查看任务配置" :visible.sync="viewFormVisible" :close-on-click-modal="false">
      <el-form :model="viewForm" label-width="100px"  ref="viewForm">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="viewForm.name" auto-complete="off" disabled></el-input>
        </el-form-item>
        <el-form-item label="任务集群" prop="group">
      <el-select v-model="viewForm.group" placeholder="请选择任务集群" style="width:100%;" disabled>
          <el-option
            v-for="item in groups"
            :key="item.name"
            :label="item.remark"
            :value="item.name">
          </el-option>
        </el-select>
        </el-form-item>
        <el-form-item label="Cron表达式" prop="cron" >
             <el-input v-model="viewForm.cron" auto-complete="off" placeholder="请输入Cron表达式" disabled></el-input>
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-select v-model="viewForm.state" placeholder="请选择任务状态" style="width:100%;" disabled>
         <el-option label="开启" :value="1"></el-option>
            <el-option label="关闭" :value="2"></el-option>
        </el-select>
        </el-form-item>
         <el-form-item label="目标任务" prop="target">
            <el-input v-model="viewForm.target" placeholder="请输入目标任务" auto-complete="off" disabled></el-input>
        </el-form-item>
         <el-form-item label="任务参数" prop="params">
            <el-input v-model="viewForm.params" placeholder="请输入任务参数" auto-complete="off" disabled></el-input>
        </el-form-item>
          <el-form-item label="手机号码" prop="mobile">
            <el-input v-model="viewForm.mobile" placeholder="请输入手机号码" auto-complete="off" disabled></el-input>
          </el-form-item>
            <el-form-item label="备注" prop="remark">
            <el-input v-model="viewForm.remark" placeholder="请输入任务备注" auto-complete="off" disabled></el-input>
          </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="viewFormVisible = false" size="mini">关闭</el-button>
      </div>
    </el-dialog>



    </div>
</template>

<script>

import {findGroupList} from '@/api/group'
import {findJobList,addJob,updateJob,deleteJob,executeJob} from '@/api/job'
export default {
    
    data(){
        return {
          groups:[],
          jobs:[],
          searchForm:{
            "group":''
          },
          addLoading:false,
          addFormVisible:false,
          addForm:{
              name: '',
              group: '',
              cron:'',
              target: '',
              param:'',
              mobile:'',
              remark:''
          },
          addFormRules:{
               name: [
          { required: true, message: '请输入任务名称', trigger: 'blur' }
        ],
         group: [
          { required: true, message: '请选择任务集群', trigger: 'blur' }
        ],
         cron: [
          { required: true, message: '请输入Cron表达式', trigger: 'blur' }
        ],
         state: [
          { required: true, message: '请选择任务状态', trigger: 'blur' }
        ],
         target: [
          { required: true, message: '请输入Target', trigger: 'blur' }
        ]
          },

      editFormVisible: false, // 新增界面是否显示
      editLoading: false,
      editFormRules: {
        name: [
          { required: true, message: '请输入任务名称', trigger: 'blur' }
        ],
         group: [
          { required: true, message: '请选择任务集群', trigger: 'blur' }
        ],
         cron: [
          { required: true, message: '请输入Cron表达式', trigger: 'blur' }
        ],
         state: [
          { required: true, message: '请选择任务状态', trigger: 'blur' }
        ],
         target: [
          { required: true, message: '请输入目标任务', trigger: 'blur' }
        ]
      },
      // 新增界面数据
      editForm: {
        id:'',
        name: '',
        group: '',
        cron:'',
        target: '',
        params:'',
        mobile:'',
        remark:''
      },
      viewForm:{
      },
      viewFormVisible:false,
        }
    },
    created(){

        this.loadGroupList()
        this.searchJobList()
    },
    methods:{

    /**
     * 删除任务配置
     */
    handleDelete(index, row){
        this.$confirm('确认要删除此任务配置信息吗？', '友情提示', {}).then(() => {
            this.addLoading = true
           
            let para = {id:row.id,group:row.group}
            deleteJob(para).then((res) => {
              this.addLoading = false
              // NProgress.done();
              if(res.code ==0){
                  this.$message({
                message: "删除成功",
                type: 'success'
              })
              this.searchJobList()
              }else{
                this.$message({
                  message: res.message,
                  type: 'error'
                })
              }
             
            })
          })


          
      },
      /**
       * 立即执行任务
       */
    handleExecute(index,row){
      this.$confirm('确认要立即此任务吗？', '友情提示', {}).then(() => {
            let para = {id:row.id,"group":row.group}
            executeJob(para).then((res) => {
                this.$message({
                      message: res.data,
                      type: 'success'
                    })
            })
       })
    },

    /**
     * 查看任务配置详情
     */
    handleView(index, row){
        this.viewForm = Object.assign({}, row)
        this.viewFormVisible =true
          
      },
      /**
       * 打开编辑弹框
       */
      handleEdit(index, row){
          this.editForm = Object.assign({}, row)
          this.editFormVisible =true
       
      },
    /**
     * 修改任务配置
     */
    editSubmit(){
      this.$refs.editForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认要修改此任务配置信息吗？', '友情提示', {}).then(() => {
            this.editLoading = true
            let para = Object.assign({}, this.editForm)
            updateJob(para).then((res) => {
              this.editLoading = false
              this.$message({
                    message: "修改成功",
                    type: 'success'
                  })
              // reset the form 
              this.$refs['editForm'].resetFields()
              this.editFormVisible = false
              this.searchJobList()
            }).catch((error)=>{
                 this.editLoading = false
            })
          })
        }
      })
   },

      /**
       * 新增任务配置
       */
      addSubmit(){
            this.$refs.addForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认要新增此任务配置信息吗？', '友情提示', {}).then(() => {
            this.addLoading = true
            let param = Object.assign({}, this.addForm)
            addJob(param).then((res) => {
              this.addLoading = false
              this.$message({
                    message: "添加成功",
                    type: 'success'
                  })
              // reset the form 
              this.$refs['addForm'].resetFields()
              this.addFormVisible = false
              this.searchJobList()
            }).catch((error)=>{
                 this.addLoading = false
            })
          })
        }
      })

      },

      /**
       * 任务配置列表搜索
       */
      onSearchSubmit(){
        this.searchJobList()
      },
      /**
       * 搜索任务配置列表
       */
      searchJobList(){
        let param =Object.assign({},this.searchForm)
        findJobList(param).then(resp =>{
            this.jobs = resp.data

        })
      },

      /**
       * 查询任务集群列表
       */
      loadGroupList(){
        findGroupList({}).then(resp =>{
          this.groups = resp.data
        })
      }
    }
}
</script>