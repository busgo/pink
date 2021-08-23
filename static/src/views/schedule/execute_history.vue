<template>
    <div class="app-container">
        <el-card class="box-card">
  <div slot="header" class="clearfix">
    <span>执行历史快照列表</span>
  </div>

  <el-form :inline="true" :model="searchForm" class="demo-form-inline" size="small">
  <!-- <el-form-item label="任务集群">
    <el-select v-model="searchForm.group" style="width:100%"   placeholder="请选择集群">
      <el-option v-for="item in groups" :key="item.name" :label="item.remark" :value="item.name"></el-option>
    </el-select>
  </el-form-item>

   <el-form-item label="IP">
    <el-input type="text" v-model="searchForm.ip" placeholder="请输入ip"></el-input>
  </el-form-item> -->

   <el-form-item label="ID">
    <el-input type="text" v-model="searchForm.id" placeholder="请输入执行快照id"></el-input>
  </el-form-item>

  <el-form-item>
    <el-button type="primary" @click="onSearchSubmit" icon="el-icon-search" round>查 询</el-button>
  </el-form-item>
</el-form>
  <el-table
    :data="snapshots"
    style="width: 100%">
      <el-table-column
      fixed
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
      width="200">
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
      prop="ip"
      label="IP"
      width="150">
    </el-table-column>
     <el-table-column
      prop="state"
      label="状态"
      width="100">
    <template slot-scope="scope">

      <el-tag v-if="scope.row.state==0" type="primary" size="mini">待执行</el-tag>
      <el-tag v-if="scope.row.state==1" type="warning" size="mini">执行中</el-tag>
      <el-tag v-if="scope.row.state==2" type="success" size="mini">成功</el-tag>
      <el-tag v-if="scope.row.state==3" type="danger" size="mini">失败</el-tag>

    </template>

    </el-table-column>

   
    <el-table-column
      prop="times"
      label="耗时"
      width="100">
         <template slot-scope="scope">
           {{scope.row.times}}秒
         </template>
    </el-table-column>
    <el-table-column
      prop="start_time"
      label="开始时间"
      width="200">
    </el-table-column>
       <el-table-column
      prop="end_time"
      label="结束时间"
      width="200">
    </el-table-column>
  
     <el-table-column
      prop="schedule_time"
      label="调度时间"
      width="200">
    </el-table-column>
     <el-table-column
      prop="cron"
      label="Cron 表达式"
      width="150">
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
      fixed="right"
      label="操作" width="80">
      <template slot-scope="scope">
            <el-button type="danger" size="mini" @click="handleDelete(scope.$index, scope.row)" round>删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</el-card>
    </div>
</template>

<script>
import {findGroupList} from '@/api/group'
import {findExecuteHistorySnapshots,deleteExecuteHistorySnapshots}from '@/api/schedule'
export default {
    
    data(){
        return {
          groups:[],
            searchForm:{

            },
            snapshots:[],
            loading:false
          
        }
    },
    created(){
      this.searchExecuteHistorySnapshots()
     // this.searchGroupList()
    },
    methods:{
 /**
       * 任务配置列表搜索
       */
      onSearchSubmit(){
        this.searchExecuteHistorySnapshots()
      },
      /**
       * 查询任务集群列表
       */
      searchGroupList(){
        findGroupList({}).then(resp =>{
          this.groups = resp.data
        })
      },
        /**
         * 删除任务执行快照
         */
        handleDelete(index, row){
                this.$confirm('确认要删除此任务执行历史快照信息吗？', '友情提示', {}).then(() => {
            this.addLoading = true
           
            let para = {id:row.id}
            deleteExecuteHistorySnapshots(para).then((res) => {
              this.addLoading = false
              // NProgress.done();
              if(res.code ==0){
                  this.$message({
                message: "删除成功",
                type: 'success'
              })
              this.searchExecuteHistorySnapshots()
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
         * 搜索任务执行快照列表
         */
        searchExecuteHistorySnapshots(){
           let param =Object.assign({},this.searchForm)
          findExecuteHistorySnapshots(param).then(resp =>{
            this.snapshots = resp.data
        })
        }
    }
}
</script>