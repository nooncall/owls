<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="模糊查询">
          <el-input v-model="searchInfo.key" placeholder="模糊查询"/>
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button size="small" icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="small" type="primary" icon="plus" @click="openDialog('add')">新增</el-button>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" row-key="id" @selection-change="handleSelectionChange">
        <el-table-column align="left" label="ID" min-width="90" prop="id" sortable="custom"/>
        <el-table-column align="left" label="类型" min-width="90" prop="data_type" sortable="custom"/>
        <el-table-column align="left" label="用户" min-width="90" prop="username" sortable="custom"/>
        <el-table-column align="left" label="集群" min-width="120" prop="cluster" sortable="custom"/>
        <el-table-column align="left" label="DB" min-width="90" prop="db" sortable="custom"/>
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
                icon="delete"
                size="small"
                type="text"
                @click="deleteAuth(scope.row)"
            >删除
            </el-button>
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
  </div>
</template>

<script lang="ts">
export default {
  name: 'Api',
}
</script>

<script lang="ts" setup>
import {
  listAuth,
  createAuth,
  delAuth,
} from '@/api/auth/auth'
import {toSQLLine} from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import {onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import moment from 'moment'

const methodFiletr = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const newLineFormatter = (row, column) => {
  return row.sql_content.replaceAll("\n", "<br>")
}

const apis = ref([])
const form = ref({
  name: '申请权限',
  db_name: '',
  task_type: '',
  remark: '',
  sql_content: ''
})
const taskType = 'auth'

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
  name: [{required: true, message: '请输入名称', trigger: 'blur'}],
  addr: [
    {required: true, message: '请输入地址', trigger: 'blur'}
  ],
  pwd: [
    {required: true, message: '请输入密码', trigger: 'blur'}
  ],
  user: [
    {required: true, message: '请输入用户名', trigger: 'blur'}
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

const dateFormatter = (row, column) => {
  return moment(row.ct * 1000).format('YYYY-MM-DD HH:mm');
}

// 排序
const sortChange = ({prop, order}) => {
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
const getTableData = async () => {
  const table = await listAuth({page: page.value, pageSize: pageSize.value, ...searchInfo.value}, taskType)
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

// 弹窗相关
const apiForm = ref(null)
const initForm = () => {
  apiForm.value.resetFields()
  form.value = {
    name: '申请权限',
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  }
}

const editTaskFunc = async (row) => {
  form.value = row
  openDialog('edit')
}

const deleteAuth = async (row) => {
  ElMessageBox.confirm('确定删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
      .then(async () => {
        const res = await delAuth(row.id)
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '操作成功!'
          })
          if (tableData.value.length === 1 && page.value > 1) {
            page.value--
          }
          getTableData()
        }
      })
}

const state = ref('')

let timeout: NodeJS.Timeout
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

.el-table .cell {
  white-space: pre-line;
}

</style>
