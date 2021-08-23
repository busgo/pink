import request from '@/utils/request'

/**
 * 查询任务调度快照列表
 * @param {*} data 
 * @returns 
 */
export function findScheduleSnapshots(data) {
    return request({
      url: '/schedule/snapshots',
      method: 'post',
      data,
    })
  }

/**
 * 删除调度快照
 * @param {*} data 
 * @returns 
 */
export function deleteScheduleSnapshot(data){
  return request({
    url:'/schedule/snapshot/delete',
    method:'post',
    data,
  })
}

/**
 * 查询执行快照列表
 * @param {*} data 
 * @returns 
 */
export function findExecuteSnapshots(data){
  return request({
    url:'/execute/snapshots',
    method:'post',
    data,
  })
}

/**
 * 查询执行历史快照列表
 * @param {*} data 
 * @returns 
 */
 export function findExecuteHistorySnapshots(data){
  return request({
    url:'/execute/history/snapshots',
    method:'post',
    data,
  })
}

/**
 * 删除执行历史快照
 * @param {*} data 
 * @returns 
 */
 export function deleteExecuteHistorySnapshots(data){
  return request({
    url:'/execute/history/snapshot/delete',
    method:'post',
    data,
  })
}