<template>
    <div class="app-container">
        <el-card class="box-card">
  <div slot="header" class="clearfix">
    <span>执行历史快照列表</span>
  </div>

  <el-form :inline="true" :model="searchForm" label-position="right" label-width="80px" class="demo-form-inline" size="small">
    <el-row>
<el-col :span="6">
      <el-form-item label="任务Id">
    <el-input type="text" v-model="searchForm.job_id" placeholder="请输入执行快照id"></el-input>
  </el-form-item>
  </el-col>
     <el-col :span="6">
<el-form-item label="任务集群">
    <el-select v-model="searchForm.group" style="width:100%" :clearable="true"   placeholder="请选择任务集群">
      <el-option v-for="item in groups" :key="item.name" :label="item.remark" :value="item.name"></el-option>
    </el-select>
  </el-form-item>
     </el-col>
      <el-col :span="12">
  <el-form-item label="调度时间">

    <el-date-picker
      v-model="searchForm.schedule_start_time"
      type="datetime"
      format="yyyy-MM-dd HH:mm:ss"
      value-format="yyyy-MM-dd HH:mm:ss"
      placeholder="选择调度开始时间">
    </el-date-picker> 至
     <el-date-picker
      v-model="searchForm.schedule_end_time"
      type="datetime"
      format="yyyy-MM-dd HH:mm:ss"
      value-format="yyyy-MM-dd HH:mm:ss"
      placeholder="选择调度结束时间">
    </el-date-picker>


  </el-form-item>
      </el-col>
    </el-row>

    <el-row>
       <el-col :span="6">
<el-form-item label="快照ID">
    <el-input type="text" v-model="searchForm.snapshot_id" placeholder="请输入快照Id"></el-input>
  </el-form-item>
       </el-col>
 <el-col :span="6">
<el-form-item label="IP">
    <el-input type="text" v-model="searchForm.ip" placeholder="请输入ip"></el-input>
  </el-form-item>
       </el-col>
        <el-col :span="6">
          <el-form-item label="执行状态">
    <el-select v-model="searchForm.state" :clearable="true" style="width:100%"   placeholder="请选执行状态">
      <el-option label="成功" value="2"></el-option>
      <el-option label="失败" value="3"></el-option>
    </el-select>
  </el-form-item>

        </el-col>

          <el-col :span="6">

            <el-form-item>
    <el-button type="primary" @click="onSearchSubmit" icon="el-icon-search" round>查 询</el-button>
  </el-form-item>
          </el-col>

    </el-row>
</el-form>

  <el-table
    :data="snapshots"
    style="width: 100%">
      <el-table-column
      fixed
      prop="id"
      label="Id"
      width="80"
      >
    </el-table-column>
    <el-table-column
      prop="job_id"
      label="任务Id"
      width="200"
      >
    </el-table-column>
    <el-table-column
          prop="snapshot_id"
          label="快照Id"
          width="200">
        </el-table-column>

    <el-table-column
      prop="job_name"
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
 
   <div class="block" style="float:right; margin-top:10px;">
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="searchForm.page_no"
      :page-sizes="[10, 20, 30, 50]"
      :page-size="searchForm.page_size"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total">
    </el-pagination>
  </div>

</el-card>
    </div>
</template>

<script>
import {findGroupList} from '@/api/group'
import {findExecuteHistorySnapshots,deleteExecuteHistorySnapshots}from '@/api/schedule'
export default {

    data(){
        return {
          value1:"",
          groups:[],
          total:100,
            searchForm:{
                schedule_start_time:'',
                schedule_start_time:'',
            },
            snapshots:[],
            loading:false

        }
    },
    created(){
      this.searchExecuteHistorySnapshots()
       this.searchGroupList()
    },
    methods:{

      handleSizeChange(pageSize){
        this.searchForm.page_size =pageSize
        this.searchForm.page_no =1
        this.searchExecuteHistorySnapshots()
      },
      handleCurrentChange(pageNo){
        this.searchForm.page_no =pageNo
        this.searchExecuteHistorySnapshots()

      },
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
