<template>
    <div class="app-container">
        <el-card class="box-card">
  <div slot="header" class="clearfix">
    <span>任务调度节点</span>
  </div>
  <el-table
    :data="nodes"
    style="width: 100%">
      <el-table-column
      prop="id"
      label="Id"
      width="400"
      >
    </el-table-column>
    <el-table-column
      prop="state"
      label="状态">
      <template slot-scope="scope">
         <el-tag  v-if="scope.row.state ==1" type="primary" size="medium">Follower</el-tag>
          <el-tag  v-if="scope.row.state ==2" type="danger" size="medium">Leader</el-tag>
      </template>
    </el-table-column>
  
  </el-table>
</el-card>
    </div>
</template>

<script>

import {findNodeList}from '@/api/node'
export default {
    
    data(){
        return {
            nodes:[]
          
        }
    },
    created(){
      this.searchNodeList()
    },
    methods:{
        
        searchNodeList(){
          findNodeList({}).then(resp =>{
            this.nodes = resp.data
        })
        }
    }
}
</script>