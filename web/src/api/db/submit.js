import service from '@/utils/request'

export const listTask = (data) => {
  return service({
    url: '/db/task/list',
    method: 'post',
    data
  })
}

export const deleteTask = (data) => {
  return service({
    url: '/db/task',
    method: 'delete',
    data
  })
}

export const updateTask = (data) => {
  return service({
    url: '/db/task',
    method: 'put',
    data
  })
}

export const createTask = (data) => {
  return service({
    url: '/db/task',
    method: 'post',
    data
  })
}
