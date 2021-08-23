import request from '@/utils/request'

/**
 *  查询任务列表
 * @param {参数} data 
 * @returns 
 */
export function findJobList(data){
    return request({
        url:'/job/list',
        method:'post',
        data,
    })
}

/**
 * 添加任务配置
 * @param {*} data 
 */
export function addJob(data){
    return request({
            url:'/job/add',
            method:'post',
            data,
    })
}

/**
 * 修改任务配置
 * @param {*} data 
 */
 export function updateJob(data){
    return request({
            url:'/job/update',
            method:'post',
            data,
    })
}


/**
 * 删除任务配置
 * @param {*} data 
 */
 export function deleteJob(data){
    return request({
            url:'/job/delete',
            method:'post',
            data,
    })
}


/**
 * 立即执行任务
 * @param {*} data 
 */
 export function executeJob(data){
    return request({
            url:'/job/execute',
            method:'post',
            data,
    })
}