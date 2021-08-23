<template>
    <div class="app-container">
        <el-card class="box-card">
  <div slot="header" class="clearfix">
    <span>调度执行计划</span>
  </div>
  <el-table
    :data="plans"
    style="width: 100%">
      <el-table-column
      fixed
      prop="id"
      label="Id"
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
      prop="next_time"
      label="下次调度时间"
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
      label="操作" width="150">
      <template slot-scope="scope">
        <el-button type="success" size="mini" @click="handleViewClients(scope.$index, scope.row)" round>任务节点</el-button>
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
export default {
    
    data(){
        return {
            plans:[],
            clients:[],
            loading:false,
            viewClientsFormVisible:false,
          
        }
    },
    created(){
      this.searchSchedulePlanList()
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
         * 搜索调度计划列表
         */
        searchSchedulePlanList(){
          findSchedulePlanList({}).then(resp =>{
            this.plans = resp.data
        })
        }
    }
}
</script>