import service from '@/utils/request'

export const listTask = (data, subType) => {
  return service({
    url: '/task/list?type=' + subType,
    method: 'post',
    data
  })
}

export const listReviewTask = (data, subType) => {
  return service({
    url: '/task/review?type=' + subType,
    method: 'post',
    data
  })
}

export const listHistoryTask = (data, subType) => {
  return service({
    url: '/task/history?type=' + subType,
    method: 'post',
    data
  })
}
export const getTask = (data, subType) => {
  return service({
    url: '/task?id=' + data + '&?type=' + subType,
    method: 'get',
  })
}

export const updateTask = (data) => {
  return service({
    url: '/task',
    method: 'put',
    data
  })
}

export const createTask = (data) => {
  return service({
    url: '/task',
    method: 'post',
    data
  })
}
