import service from '@/utils/request'

export const readRedisData = (data) => {
  return service({
    url: '/redis/read',
    method: 'post',
    data
  })
}


export const readData = (data) => {
  return service({
    url: '/db/read',
    method: 'post',
    data
  })
}

export const getTableInfo = (data) => {
  return service({
    url: '/db/table',
    method: 'post',
    data
  })
}
