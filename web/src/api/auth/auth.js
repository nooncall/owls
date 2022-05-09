import service from '@/utils/request'

export const listAuth = (data, subType) => {
  return service({
    url: '/auth/list',
    method: 'post',
    data
  })
}

export const delAuth = (data) => {
  return service({
    url: '/auth?id=' + data,
    method: 'delete',
  })
}

export const createAuth = (data) => {
  return service({
    url: '/auth',
    method: 'post',
    data
  })
}
