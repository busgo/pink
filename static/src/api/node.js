import request from '@/utils/request'

/**
 * 查询当前调度节点列表
 * @param {参数} data 
 * @returns 
 */
export function findNodeList(data) {
  return request({
    url: '/node/list',
    method: 'post',
    data,
  })
}