import service from '@/utils/request'

export const listCluster = (data) => {
  return service({
    url: '/db/cluster/list',
    method: 'post',
    data
  })
}

export const listClusterName = (filter, ctype) => {
  return service({
    url: '/db/cluster/name/list?filter=' + filter + '&c_type=' + ctype,
    method: 'get',
  })
}

export const listDatabase = (cluster, filter) => {
  return service({
    url: '/db/cluster/db/list?filter=' + filter + '&cluster=' + cluster,
    method: 'get',
  })
}

export const listTable = (cluster, db) => {
  return service({
    url: '/db/cluster/table/list?cluster=' + cluster + '&db=' + db,
    method: 'get'
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
