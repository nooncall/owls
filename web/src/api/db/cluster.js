import service from '@/utils/request'

export const listCluster = (data) => {
  return service({
    url: '/db/cluster/list',
    method: 'post',
    data
  })
}

export const deleteCluster = (data) => {
  return service({
    url: '/db/cluster?id=' + data.id,
    method: 'delete',
    data
  })
}

export const updateCluster = (data) => {
  return service({
    url: '/db/cluster',
    method: 'put',
    data
  })
}

export const createCluster = (data) => {
  return service({
    url: '/db/cluster',
    method: 'post',
    data
  })
}
