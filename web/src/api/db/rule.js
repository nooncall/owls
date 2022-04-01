import service from '@/utils/request'

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
