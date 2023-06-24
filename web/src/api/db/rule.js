import service from '@/utils/request'

export const listRedisRule = (data) => {
  return service({
    url: '/redis/rule/list',
    method: 'get',
    data
  })
}

export const listRule = (data) => {
  return service({
    url: '/db/rule/list',
    method: 'post',
    data
  })
}

export const updateRuleStatus = (data) => {
  return service({
    url: '/db/rule/status',
    method: 'put',
    data
  })
}
