import request from '@/utils/request'

/**
 * 查询当前集群列表
 * @param {参数} data 
 * @returns 
 */
export function findGroupList(data) {
  return request({
    url: '/group/list',
    method: 'post',
    data,
  })
}


/**
 * 查询当前集群详情列表
 * @param {参数} data 
 * @returns 
 */
 export function findGroupDetailsList(data) {
  return request({
    url: '/group/details/list',
    method: 'post',
    data,
  })
}

/**
 * 添加任务集群
 * @param {*} data 
 * @returns 
 */
export function addGroup(data){
  return request({
    url:'/group/add',
    method:'post',
    data,
  })
}