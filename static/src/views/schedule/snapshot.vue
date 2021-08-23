<template>
    <div class="app-container">
      <el-alert
      style="margin-bottom:10px;"
    title="展示未分配的任务调度快照"
    type="warning"
    :closable="false"
    show-icon>
  </el-alert>
  <el-card class="box-card">
  <div slot="header" class="clearfix">
    <span>调度快照列表</span>
  </div>
  <el-table
    :data="snapshots"
    style="width: 100%">
      <el-table-column
      prop="id"
      label="Id"
      width="200"
      >
    </el-table-column>
    <el-table-column
      prop="job_id"
      label="任务ID"
      width="200"
      >
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
      <template slot-scope="scope">
            <el-tag  type="primary" size="medium" style="panding-right:3px;">{{scope.row.group}}</el-tag> 
      </template>
    </el-table-column>
    <el-table-column
      prop="cron"
      label="Cron 表达式"
      width="150">
    </el-table-column>
 <el-table-column
      prop="before_time"
      label="上次调度时间"
      width="200">
    </el-table-column>
     <el-table-column
      prop="schedule_time"
      label="调度时间"
      width="200">
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
      fixed="right"
      label="操作" width="350">
      <template slot-scope="scope">
           <el-button type="success" size="mini" @click="handleViewClients(scope.$index, scope.row)" round>任务节点</el-button>
            <el-button type="danger" size="mini" @click="handleDelete(scope.$index, scope.row)" round>删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</el-card>

  <!--查看界面-->
    <el-dialog title="执行任务Client列表" :visible.sync="viewClientsFormVisible" :close-on-click-modal="false">
      
       <!--列表-->
      <el-table :data="clients" highlight-current-row v-loading="loading"  style="width: 100%;">
        
        <el-table-column prop="ip" label="client" width="150">
        </el-table-column>
        <el-table-column prop="group" label="集群" width="150">
        </el-table-column>
        <el-table-column prop="path" label="path">
        </el-table-column>
      </el-table>

      <div slot="footer" class="dialog-footer">
        <el-button @click.native="viewClientsFormVisible = false" size="mini">关闭</el-button>
      </div>
    </el-dialog>



    </div>
</template>

<script>
import {findSchedulePlanList,findSchedulePlanClients}from '@/api/plan'
import {findScheduleSnapshots,deleteScheduleSnapshot}from '@/api/schedule'
export default {
    
    data(){
        return {
            clients:[],
            snapshots:[],
            loading:false,
            viewClientsFormVisible:false,
          
        }
    },
    created(){
      this.searchScheduleSnapshots()
    },
    methods:{

            handleViewClients(index, row){
             this.viewClientsFormVisible =true
                let param ={
                    group:row.group
                }
                findSchedulePlanClients(param).then((res) => {
                this.loading = false
                this.clients = res.data
            })
          },
        /**
         * 删除任务调度快照
         */
        handleDelete(index, row){
                this.$confirm('确认要删除此调度快照信息吗？', '友情提示', {}).then(() => {
            this.addLoading = true
           
            let para = {id:row.id}
            deleteScheduleSnapshot(para).then((res) => {
              this.addLoading = false
              // NProgress.done();
              if(res.code ==0){
                  this.$message({
                message: "删除成功",
                type: 'success'
              })
              this.searchScheduleSnapshots()
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
         * 搜索任务调度快照列表
         */
        searchScheduleSnapshots(){
          findScheduleSnapshots({}).then(resp =>{
            this.snapshots = resp.data
        })
        }
    }
}
</script>