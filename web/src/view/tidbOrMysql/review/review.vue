<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="模糊查询">
          <el-input v-model="searchInfo.key" placeholder="模糊查询" />
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button size="small" icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" @sort-change="sortChange" row-key="id" @selection-change="handleSelectionChange">
        <el-table-column type="expand">
          <template #default="scope">
            <el-table :data="scope.row.exec_items" style="width: calc(100% - 47px)" class="two-list">
              <el-table-column prop="id" label="序号"></el-table-column>
              <el-table-column prop="cluster_name" label="集群"></el-table-column>
              <el-table-column prop="db_name" label="库名"></el-table-column>
              <el-table-column prop="task_type" label="类型"></el-table-column>
              <el-table-column prop="affect_rows" label="影响行数"></el-table-column>
              <el-table-column prop="status" label="状态"></el-table-column>
              <el-table-column prop="exec_info" label="执行信息"></el-table-column>
              <el-table-column class="cell" prop="cat_id" width="600"  label="SQL语句">
                <template class="cell" style="white-space: pre-line;" #default="scope">
                  <code>
                    <span v-html="newLineFormatter(scope.row, '')"></span>
                  </code>
                </template>
              </el-table-column>
              <el-table-column prop="remark" label="备注" width="200"></el-table-column>
            </el-table>
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="150" prop="id" sortable="custom" />
        <el-table-column align="left" label="任务名" min-width="150" prop="name" sortable="custom" />
        <el-table-column align="left" label="状态" min-width="150" prop="status" sortable="custom" />
        <el-table-column align="left" label="创建者" min-width="150" prop="creator" sortable="custom" />
        <el-table-column align="left" label="创建时间" min-width="150" prop="ct" :formatter="dateFormatter" sortable="custom" />
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
                icon="delete"
                size="small"
                type="text"
                @click="rejectFunc(scope.row)"
            >驳回</el-button>
            <el-button
                icon="delete"
                size="small"
                type="text"
                @click="toTarget('exec')"
            >审核</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 20, 30, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="rejectDialogFormVisible" :before-close="closeDialog" :title="驳回">
      <el-form ref="apiForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="驳回原因" prop="description">
          <el-input
              v-model="form.reject_content"
              :autosize="{ minRows: 3, maxRows: 5000 }"
              type="textarea"
              placeholder="Please input"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterEditDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
export default {
  name: 'Api',
}
</script>

<script lang="ts" setup>
import {
  listTask,
  createTask,
  updateTask,
  cancelTask,
} from '@/api/db/submit'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import moment from 'moment'
import {
  listCluster,
  listDatabase,
} from '@/api/db/cluster'
import { useRouter } from 'vue-router'


const methodFiletr = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const router = useRouter()
const toTarget = (name) => {
  console.log("name is : ", name)
  router.push(name)
}

const newLineFormatter = (row, column) =>{
  return row.sql_content.replaceAll("\n", "<br>")
}

const apis = ref([])
const form = ref({
  cluster_name: '',
  db_name:'',
  task_type: '',
  remark: '',
  sql_content: ''
})
const methodOptions = ref([
  {
    value: 'POST',
    label: '创建',
    type: 'success'
  },
  {
    value: 'GET',
    label: '查看',
    type: ''
  },
  {
    value: 'PUT',
    label: '更新',
    type: 'warning'
  },
  {
    value: 'DELETE',
    label: '删除',
    type: 'danger'
  }
])

const type = ref('')
const rules = ref({
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  addr: [
    { required: true, message: '请输入地址', trigger: 'blur' }
  ],
  pwd: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  user: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const onReset = () => {
  searchInfo.value = {}
}
// 搜索

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const dateFormatter = (row, column) =>{
  return moment(row.ct *1000).format('YYYY-MM-DD HH:mm');
}

// 排序
const sortChange = ({ prop, order }) => {
  if (prop) {
    if (prop === 'ID') {
      prop = 'id'
    }
    searchInfo.value.orderKey = toSQLLine(prop)
    searchInfo.value.desc = order === 'descending'
  }
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await listTask({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// 批量操作
const handleSelectionChange = (val) => {
  apis.value = val
}

const deleteVisible = ref(false)
const onDelete = async() => {
  const ids = apis.value.map(item => item.ID)
  const res = await deleteApisByIds({ ids })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    deleteVisible.value = false
    getTableData()
  }
}

// 弹窗相关
const apiForm = ref(null)
const initForm = () => {
  apiForm.value.resetFields()
  form.value = {
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  }
}

const dialogTitle = ref('新增')
const handleId = ref('');
const rejectDialogFormVisible = ref(false)
const openDialog = (id) => {
  rejectDialogFormVisible.value = true
  handleId.value = id
}

const closeDialog = () => {
  initForm()
  rejectDialogFormVisible.value = false
}

const editTaskFunc = async(row) => {
  form.value = row
  openDialog('edit')
}

const enterEditDialog = async() => {
  apiForm.value.validate(async valid => {
    let params = {
      id: handleId.value,
      reject_content: form.value.reject_content,
      action: "reject"
    }
    const res = await updateTask(params)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '操作成功',
        showClose: true
      })
    }
    getTableData()
    closeDialog()
  })
}


const enterDialog = async() => {
  apiForm.value.validate(async valid => {
    if (valid) {
      switch (type.value) {
        case 'add':
        {
          //todo, refactor
          let paramas = {
            name: form.value.name,
            sub_tasks:[
              {
                cluster_name: form.value.cluster_name,
                db_name: form.value.db_name,
                task_type: form.value.task_type,
                exec_items:[
                  {
                    remark: form.value.remark,
                    sql_content: form.value.sql_content,
                  }
                ]
              }
            ]
          }
          const res = await createTask(paramas)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '添加成功',
              showClose: true
            })
          }
          getTableData()
          closeDialog()
        }

          break
        case 'edit':
        {
          const res = await updateTask(form.value)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '编辑成功',
              showClose: true
            })
          }
          getTableData()
          closeDialog()
        }
          break
        default:
          // eslint-disable-next-line no-lone-blocks
        {
          ElMessage({
            type: 'error',
            message: '未知操作',
            showClose: true
          })
        }
          break
      }
    }
  })
}

const rejectFunc = async(row) => {
  handleId.value = row.id
  rejectDialogFormVisible.value = true
}

const state = ref('')

interface clusterItem {
  name: string
  value: string
}

const clusters = ref<clusterItem[]>([])
const db = ref<clusterItem[]>([])

const loadCluster = async () => {
  const resp = await listCluster({page: 1, pageSize: 5})
  let result = []
  for (let cluster of resp.data.list){
    result.push({value: cluster.name, name: cluster.name})
  }

  return result
}

const loadDB = async (cluster) => {
  const resp = await listDatabase(cluster)
  let result = []
  for (let db of resp.data){
    result.push({value: db, name: db})
  }

  return result
}

let timeout: NodeJS.Timeout
const querySearchAsync = (queryString: string, cb: (arg: any) => void) => {
  const results = queryString
      ? clusters.value.filter(createFilter(queryString))
      : clusters.value

  clearTimeout(timeout)
  timeout = setTimeout(() => {
    cb(results)
  }, 3000 * Math.random())
}

const querySearchDBAsync = async(queryString: string, cb: (arg: any) => void) => {
  const results = queryString
      ? db.value.filter(createFilter(queryString))
      : db.value

  clearTimeout(timeout)
  timeout = setTimeout(() => {
    cb(results)
  }, 3000 * Math.random())
}

const createFilter = (queryString: string) => {
  return (restaurant: clusterItem) => {
    return (
        restaurant.name.toLowerCase().indexOf(queryString.toLowerCase()) === 0
    )
  }
}

const handleSelect = async (item: clusterItem) => {
  db.value = await loadDB(item.value)
}

onMounted(async () => {
  clusters.value = await loadCluster()
})

</script>


<style scoped lang="scss">
.button-box {
  padding: 10px 20px;
  .el-button {
    float: right;
  }
}
.warning {
  color: #dc143c;
}

.el-table .cell{
  white-space:pre-line;
}

</style>
