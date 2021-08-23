import request from '@/utils/request'

/**
 * 查询调度计划列表
 * @param {*} data 
 * @returns 
 */
export function findSchedulePlanList(data){
    return request({
        url:'/schedule/plan/list',
        method:'post',
        data,
    })
}

/**
 * 查询调度计划客户端列表
 * @param {*} data 
 * @returns 
 */
export function findSchedulePlanClients(data){
    return request({
        url:'/schedule/plan/clients',
        method:'post',
        data,
    })
}